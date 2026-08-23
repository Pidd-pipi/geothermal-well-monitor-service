package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func alertsVerifyDeps() *AppDeps {
	return &AppDeps{
		wells:     NewWellService(NewWellStore()),
		readings:  newReadingService(newReadingStore(5000)),
		alerts:    newAlertService(newAlertStore()),
		tasks:     newTaskService(newTaskStore()),
		plans:     newPlanService(newPlanStore()),
		ops:       newOpsService(nil),
		backfill:  newBackfillRunner(newReadingStore(5000)),
		operators: newOperatorStore(defaultOperators()),
		sessions:  newSessionStore(),
	}
}

func TestAlertsListSharesNoBacking(t *testing.T) {
	store := newAlertStore()
	store.Add(Alert{ID: "alt-1", WellID: "well-07", Kind: AlertPressureHigh, Status: AlertOpen, At: time.Now()})
	store.Add(Alert{ID: "alt-2", WellID: "well-07", Kind: AlertTempHigh, Status: AlertAcked, At: time.Now()})
	all := store.List("")
	if len(all) != 2 {
		t.Fatalf("all len=%d", len(all))
	}
	all[0].ID = "mutated"
	again := store.List("")
	if again[0].ID == "mutated" {
		t.Fatalf("list result shares store backing array")
	}
}

func TestAlertsEvaluateDeduplicatesOpen(t *testing.T) {
	svc := newAlertService(newAlertStore())
	ctx := context.Background()
	reading := Reading{WellID: "well-07", PressureKPa: 1800, TemperatureC: 260}
	first, err := svc.EvaluateReading(ctx, reading)
	if err != nil {
		t.Fatal(err)
	}
	if len(first) == 0 {
		t.Fatal("no alerts created")
	}
	second, err := svc.EvaluateReading(ctx, reading)
	if err != nil {
		t.Fatal(err)
	}
	if len(second) != 0 {
		t.Fatalf("duplicate open alerts created: %+v", second)
	}
	open := svc.List(AlertOpen)
	if len(open) != len(first) {
		t.Fatalf("open count=%d want %d", len(open), len(first))
	}
}

func TestAlertsDefaultListOnlyOpen(t *testing.T) {
	svc := newAlertService(newAlertStore())
	ctx := context.Background()
	created, err := svc.EvaluateReading(ctx, Reading{WellID: "well-12", PressureKPa: 1700, TemperatureC: 240})
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 {
		t.Fatalf("created=%d", len(created))
	}
	if _, err := svc.Ack(ctx, created[0].ID, "liang"); err != nil {
		t.Fatal(err)
	}
	deps := alertsVerifyDeps()
	deps.alerts = svc
	ts := httptest.NewServer(NewRouter(deps))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/alerts")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body struct {
		Items []Alert `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	for _, a := range body.Items {
		if a.Status != AlertOpen {
			t.Fatalf("default alert list includes non-open alert %s status=%s", a.ID, a.Status)
		}
	}
}

func TestAlertAckHTTPMissing404(t *testing.T) {
	ts := httptest.NewServer(NewRouter(alertsVerifyDeps()))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/alerts/no-such-alert/ack", "application/json", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("ack missing alert status=%d (want 404)", resp.StatusCode)
	}
}

func TestAlertAckHappyPath(t *testing.T) {
	svc := newAlertService(newAlertStore())
	ctx := context.Background()
	created, err := svc.EvaluateReading(ctx, Reading{WellID: "well-07", PressureKPa: 1600, TemperatureC: 230})
	if err != nil {
		t.Fatal(err)
	}
	if len(created) == 0 {
		t.Fatal("no alert created")
	}
	acked, err := svc.Ack(ctx, created[0].ID, "liang")
	if err != nil {
		t.Fatal(err)
	}
	if acked.Status != AlertAcked {
		t.Fatalf("status=%s", acked.Status)
	}
}
