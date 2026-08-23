package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func newVerifyAppDeps(store *ReadingStore, runner *BackfillRunner) *AppDeps {
	return &AppDeps{
		wells:     NewWellService(NewWellStore()),
		readings:  newReadingService(store),
		alerts:    newAlertService(newAlertStore()),
		tasks:     newTaskService(newTaskStore()),
		plans:     newPlanService(newPlanStore()),
		ops:       newOpsService(nil),
		backfill:  runner,
		operators: newOperatorStore(defaultOperators()),
		sessions:  newSessionStore(),
	}
}

func waitBackfillIdle(t *testing.T, runner *BackfillRunner, id string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := runner.WaitIdle(ctx, id); err != nil {
		t.Fatalf("backfill %s did not drain workers: %v", id, err)
	}
}

func TestBackfillNormalCompletesAndPersists(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	payload := "well-07,812.4,142.6\nwell-07,900.0,150.0\nwell-12,700.0,130.0\n"
	job, err := runner.Start(context.Background(), payload)
	if err != nil {
		t.Fatal(err)
	}
	waitBackfillIdle(t, runner, job.ID)
	got, err := runner.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != BackfillDone {
		t.Fatalf("status=%s err=%s", got.Status, got.ErrMsg)
	}
	if got.Inserted != got.Total || got.Inserted != 3 {
		t.Fatalf("inserted=%d total=%d", got.Inserted, got.Total)
	}
	if store.Count("well-07") != 2 || store.Count("well-12") != 1 {
		t.Fatalf("store counts well-07=%d well-12=%d", store.Count("well-07"), store.Count("well-12"))
	}
}

func TestBackfillCancelledJobDrains(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	var b strings.Builder
	for i := 0; i < 200; i++ {
		b.WriteString("well-07,812.4,142.6\n")
	}
	job, err := runner.Start(context.Background(), b.String())
	if err != nil {
		t.Fatal(err)
	}
	if err := runner.Cancel(context.Background(), job.ID); err != nil {
		t.Fatal(err)
	}
	waitBackfillIdle(t, runner, job.ID)
	got, err := runner.Get(context.Background(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != BackfillFailed {
		t.Fatalf("cancelled job status=%s (want failed)", got.Status)
	}
}

func TestBackfillConcurrentCancellationRace(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var b strings.Builder
			for j := 0; j < 100; j++ {
				b.WriteString("well-12,700.0,130.0\n")
			}
			job, err := runner.Start(context.Background(), b.String())
			if err != nil {
				t.Error(err)
				return
			}
			<-start
			_ = runner.Cancel(context.Background(), job.ID)
			waitBackfillIdle(t, runner, job.ID)
		}()
	}
	close(start)
	wg.Wait()
}

func TestBackfillConcurrentAppendRecentRace(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	var b strings.Builder
	for i := 0; i < 200; i++ {
		b.WriteString("well-07,812.4,142.6\n")
	}
	job, err := runner.Start(context.Background(), b.String())
	if err != nil {
		t.Fatal(err)
	}
	start := make(chan struct{})
	stop := make(chan struct{})
	var readers sync.WaitGroup
	readers.Add(2)
	for i := 0; i < 2; i++ {
		go func() {
			defer readers.Done()
			<-start
			for {
				select {
				case <-stop:
					return
				default:
					_ = store.Recent("well-07", 10)
				}
			}
		}()
	}
	close(start)
	time.Sleep(150 * time.Millisecond)
	close(stop)
	readers.Wait()
	waitBackfillIdle(t, runner, job.ID)
}

func TestBackfillHTTPMissingID404(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	ts := httptest.NewServer(NewRouter(newVerifyAppDeps(store, runner)))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/backfill/no-such-job")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("missing backfill status=%d (want 404)", resp.StatusCode)
	}
}

func TestBackfillHTTPJobProgress(t *testing.T) {
	store := newReadingStore(5000)
	runner := newBackfillRunner(store)
	ts := httptest.NewServer(NewRouter(newVerifyAppDeps(store, runner)))
	defer ts.Close()
	var b strings.Builder
	for i := 0; i < 100; i++ {
		b.WriteString("well-07,812.4,142.6\n")
	}
	body := `{"payload":` + strconv.Quote(b.String()) + `}`
	resp, err := http.Post(ts.URL+"/api/backfill", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("start status=%d", resp.StatusCode)
	}
	var start struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&start); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		gr, err := http.Get(ts.URL + "/api/backfill/" + start.ID)
		if err != nil {
			t.Fatal(err)
		}
		var state struct {
			Status   string `json:"status"`
			Inserted int    `json:"inserted"`
			Total    int    `json:"total"`
		}
		_ = json.NewDecoder(gr.Body).Decode(&state)
		gr.Body.Close()
		if state.Status == string(BackfillDone) {
			if state.Inserted != state.Total || state.Inserted != 100 {
				t.Fatalf("premature done inserted=%d total=%d", state.Inserted, state.Total)
			}
			return
		}
		time.Sleep(30 * time.Millisecond)
	}
	t.Fatalf("backfill job %s never reached done", start.ID)
}
