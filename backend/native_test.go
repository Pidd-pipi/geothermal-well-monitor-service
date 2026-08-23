package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func testRouter() http.Handler {
	deps := &AppDeps{
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
	return NewRouter(deps)
}

func TestNativeWellFlow(t *testing.T) {
	router := testRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/health")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("health status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
	resp, err = http.Get(ts.URL + "/api/wells")
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("wells status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
	resp, err = http.Post(ts.URL+"/api/wells/status", "application/json", bytes.NewBufferString(`{"id":"well-07","status":"inspection"}`))
	if err != nil || resp.StatusCode != 200 {
		t.Fatalf("status update status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
}

func TestNativeReadingsFlow(t *testing.T) {
	svc := newReadingService(newReadingStore(100))
	ctx := context.Background()
	first, err := svc.Ingest(ctx, Reading{WellID: "well-07", PressureKPa: 812, TemperatureC: 142, Source: "sensor"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := svc.Ingest(ctx, Reading{WellID: "well-07", PressureKPa: 900, TemperatureC: 150})
	if err != nil {
		t.Fatal(err)
	}
	recent, err := svc.Recent(ctx, "well-07", 10)
	if err != nil || len(recent) != 2 {
		t.Fatalf("recent len=%d err=%v", len(recent), err)
	}
	if recent[0].ID != first.ID || recent[1].ID != second.ID {
		t.Fatalf("recent order wrong: %+v", recent)
	}
	stats, err := svc.Stats(ctx, "well-07", time.Now().Add(-time.Hour))
	if err != nil || stats.Count != 2 {
		t.Fatalf("stats count=%d err=%v", stats.Count, err)
	}
}

func TestNativeAlertsFlow(t *testing.T) {
	svc := newAlertService(newAlertStore())
	ctx := context.Background()
	created, err := svc.EvaluateReading(ctx, Reading{WellID: "well-07", PressureKPa: 1800, TemperatureC: 260})
	if err != nil || len(created) == 0 {
		t.Fatalf("created=%d err=%v", len(created), err)
	}
	if svc.OpenCount() != len(created) {
		t.Fatalf("open count mismatch")
	}
	acked, err := svc.Ack(ctx, created[0].ID, "liang")
	if err != nil || acked.Status != AlertAcked {
		t.Fatalf("ack err=%v status=%v", err, acked.Status)
	}
	open := svc.List(AlertOpen)
	ackedList := svc.List(AlertAcked)
	if len(open) != len(created)-1 || len(ackedList) != 1 {
		t.Fatalf("open=%d acked=%d", len(open), len(ackedList))
	}
}

func TestNativeTaskLifecycle(t *testing.T) {
	svc := newTaskService(newTaskStore())
	ctx := context.Background()
	task, err := svc.Create(ctx, Task{WellID: "well-12", Subject: "例行巡检", Owner: "liang", Priority: OpsPriorityHigh})
	if err != nil {
		t.Fatal(err)
	}
	task, err = svc.Assign(ctx, task.ID, "liang", task.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != TaskAssigned || task.Operator != "liang" {
		t.Fatalf("assign result %+v", task)
	}
	task, err = svc.Start(ctx, task.ID, "liang", task.Revision)
	if err != nil {
		t.Fatal(err)
	}
	task, err = svc.Complete(ctx, task.ID, "liang", task.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if task.Status != TaskCompleted {
		t.Fatalf("complete status=%s", task.Status)
	}
	history := svc.CompletedHistory()
	if len(history) != 1 || history[0].ID != task.ID {
		t.Fatalf("history=%d", len(history))
	}
}

func TestNativePlanFlow(t *testing.T) {
	svc := newPlanService(newPlanStore())
	ctx := context.Background()
	plan, err := svc.Create(ctx, Plan{WellID: "well-07", Subject: "月度保养", Owner: "zhao", IntervalDays: 30})
	if err != nil {
		t.Fatal(err)
	}
	due, err := svc.Due(ctx, time.Now().Add(31*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(due) != 1 || due[0].ID != plan.ID {
		t.Fatalf("due=%d", len(due))
	}
	done, err := svc.MarkDone(ctx, plan.ID, plan.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if !done.NextDueAt.After(time.Now()) {
		t.Fatalf("next due not advanced: %v", done.NextDueAt)
	}
}

func TestNativeOpsFlow(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	record, err := svc.Create(ctx, OpsRecord{Subject: "井口渗漏检查", Owner: "liang", Priority: OpsPriorityHigh, Labels: map[string]string{"site": "north"}})
	if err != nil {
		t.Fatal(err)
	}
	page, err := svc.Search(ctx, OpsQuery{Owner: "liang", Page: 1, PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != record.ID {
		t.Fatalf("search items=%d", len(page.Items))
	}
	updated, err := svc.Transition(ctx, record.ID, record.Revision, OpsStatusActive, "liang")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != OpsStatusActive {
		t.Fatalf("transition status=%s", updated.Status)
	}
}

func TestNativeAuthOperators(t *testing.T) {
	store := newOperatorStore(defaultOperators())
	op, err := store.Verify("liang", "well2026")
	if err != nil || op == nil || op.Username != "liang" {
		t.Fatalf("verify err=%v op=%v", err, op)
	}
	if _, err := store.Verify("liang", "wrong"); err == nil {
		t.Fatalf("wrong password accepted")
	}
}

func TestNativeBackfillParse(t *testing.T) {
	readings, err := parseBackfillPayload("well-07,812.4,142.6\nwell-12,700,130")
	if err != nil || len(readings) != 2 {
		t.Fatalf("parse len=%d err=%v", len(readings), err)
	}
	if _, err := parseBackfillPayload("well-07,bad,142.6"); err == nil {
		t.Fatalf("bad line accepted")
	}
	if _, err := parseBackfillPayload("  \n"); err == nil {
		t.Fatalf("empty payload accepted")
	}
}

func TestNativePages(t *testing.T) {
	router := testRouter()
	ts := httptest.NewServer(router)
	defer ts.Close()
	for _, path := range []string{"/", "/app.js"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil || resp.StatusCode != 200 {
			t.Fatalf("GET %s status=%v err=%v", path, resp.StatusCode, err)
		}
		resp.Body.Close()
	}
}
