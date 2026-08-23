package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func readingsVerifyDeps(store *ReadingStore) *AppDeps {
	return &AppDeps{
		wells:     NewWellService(NewWellStore()),
		readings:  newReadingService(store),
		alerts:    newAlertService(newAlertStore()),
		tasks:     newTaskService(newTaskStore()),
		plans:     newPlanService(newPlanStore()),
		ops:       newOpsService(nil),
		backfill:  newBackfillRunner(newReadingStore(5000)),
		operators: newOperatorStore(defaultOperators()),
		sessions:  newSessionStore(),
	}
}

func TestReadingRecentReturnsIsolatedCopy(t *testing.T) {
	store := newReadingStore(5000)
	svc := newReadingService(store)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		if _, err := svc.Ingest(ctx, Reading{WellID: "well-07", PressureKPa: 800 + float64(i), TemperatureC: 140 + float64(i)}); err != nil {
			t.Fatal(err)
		}
	}
	first, err := svc.Recent(ctx, "well-07", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) != 3 {
		t.Fatalf("first len=%d", len(first))
	}
	first[0].PressureKPa = -999
	second, err := svc.Recent(ctx, "well-07", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(second) != 3 {
		t.Fatalf("second len=%d", len(second))
	}
	if second[0].PressureKPa == -999 {
		t.Fatalf("mutation leaked into cached recent list")
	}
}

func TestReadingsConcurrentQueryAndAppend(t *testing.T) {
	store := newReadingStore(5000)
	svc := newReadingService(store)
	ctx := context.Background()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_, _ = svc.Ingest(ctx, Reading{WellID: "well-07", PressureKPa: 800 + float64(j), TemperatureC: 140})
			}
		}(i)
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_, _ = svc.Recent(ctx, "well-07", 10)
			}
		}(i)
	}
	close(start)
	wg.Wait()
	if svc.Count("well-07") != 200 {
		t.Fatalf("count=%d", svc.Count("well-07"))
	}
}

func TestReadingsConcurrentStatsQuery(t *testing.T) {
	store := newReadingStore(5000)
	svc := newReadingService(store)
	ctx := context.Background()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_, _ = svc.Ingest(ctx, Reading{WellID: "well-12", PressureKPa: 700, TemperatureC: 130})
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				_, _ = svc.Stats(ctx, "well-12", time.Now().Add(-time.Hour))
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestReadingHTTPInvalidRejected400(t *testing.T) {
	store := newReadingStore(5000)
	ts := httptest.NewServer(NewRouter(readingsVerifyDeps(store)))
	defer ts.Close()
	body := `{"well_id":"well-07","pressure_kpa":99999,"temperature_c":140}`
	resp, err := http.Post(ts.URL+"/api/readings", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid reading status=%d (want 400)", resp.StatusCode)
	}
}

func TestReadingsHTTPParallelIngestQuery(t *testing.T) {
	store := newReadingStore(5000)
	ts := httptest.NewServer(NewRouter(readingsVerifyDeps(store)))
	defer ts.Close()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			<-start
			for j := 0; j < 30; j++ {
				payload := fmt.Sprintf(`{"well_id":"well-07","pressure_kpa":%d,"temperature_c":140}`, 800+j)
				resp, err := http.Post(ts.URL+"/api/readings", "application/json", bytes.NewBufferString(payload))
				if err == nil {
					resp.Body.Close()
				}
			}
		}(i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 40; j++ {
				resp, err := http.Get(ts.URL + "/api/readings/recent?well_id=well-07&n=10")
				if err == nil {
					var out map[string]any
					_ = json.NewDecoder(resp.Body).Decode(&out)
					resp.Body.Close()
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestReadingStatsConsistent(t *testing.T) {
	store := newReadingStore(5000)
	svc := newReadingService(store)
	ctx := context.Background()
	for i := 0; i < 4; i++ {
		if _, err := svc.Ingest(ctx, Reading{WellID: "well-07", PressureKPa: 800 + float64(i), TemperatureC: 140}); err != nil {
			t.Fatal(err)
		}
	}
	stats, err := svc.Stats(ctx, "well-07", time.Now().Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if stats.Count != 4 {
		t.Fatalf("count=%d", stats.Count)
	}
	if strconv.FormatFloat(stats.AvgPressure, 'f', -1, 64) != strconv.FormatFloat(801.5, 'f', -1, 64) {
		t.Fatalf("avg=%v", stats.AvgPressure)
	}
}
