package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func plansVerifyDeps() *AppDeps {
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

func TestPlanMarkDoneAdvancesOverdueDueDate(t *testing.T) {
	fakeNow := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	clock := OpsClock{NowFunc: func() time.Time { return fakeNow }}
	store := newPlanStore()
	svc := &PlanService{store: store, clock: clock}
	ctx := context.Background()
	plan, err := svc.Create(ctx, Plan{WellID: "well-07", Subject: "月度保养", Owner: "zhao", IntervalDays: 30})
	if err != nil {
		t.Fatal(err)
	}
	fakeNow = fakeNow.AddDate(0, 0, 40)
	done, err := svc.MarkDone(ctx, plan.ID, plan.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if !done.NextDueAt.After(fakeNow) {
		t.Fatalf("overdue plan still due after completion: next_due=%v now=%v", done.NextDueAt, fakeNow)
	}
	due, err := svc.Due(ctx, fakeNow.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(due) != 0 {
		t.Fatalf("completed plan still listed as due: %+v", due)
	}
}

func TestPlanDueExcludesPaused(t *testing.T) {
	store := newPlanStore()
	svc := &PlanService{store: store, clock: newOpsClock()}
	ctx := context.Background()
	plan, err := svc.Create(ctx, Plan{WellID: "well-12", Subject: "换季保养", Owner: "liang", IntervalDays: 60})
	if err != nil {
		t.Fatal(err)
	}
	current, err := store.Get(ctx, plan.ID)
	if err != nil {
		t.Fatal(err)
	}
	current.NextDueAt = time.Now().Add(-time.Hour)
	current.Status = PlanPaused
	if _, err := store.Update(ctx, current, current.Revision); err != nil {
		t.Fatal(err)
	}
	due, err := svc.Due(ctx, time.Now())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range due {
		if p.ID == plan.ID {
			t.Fatalf("paused plan %s listed as due", p.ID)
		}
	}
}

func TestPlanDueHTTPFreshlyDoneExcluded(t *testing.T) {
	ts := httptest.NewServer(NewRouter(plansVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/plans", "application/json",
		strings.NewReader(`{"well_id":"well-07","subject":"季度保养","owner":"zhao","interval_days":30}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	var created Plan
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	createResp.Body.Close()
	doneResp, err := http.Post(ts.URL+"/api/plans/"+created.ID+"/done", "application/json", strings.NewReader(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if doneResp.StatusCode != http.StatusOK {
		t.Fatalf("done status=%d", doneResp.StatusCode)
	}
	doneResp.Body.Close()
	dueResp, err := http.Get(ts.URL + "/api/plans?due=true")
	if err != nil {
		t.Fatal(err)
	}
	defer dueResp.Body.Close()
	var body struct {
		Items []Plan `json:"items"`
	}
	if err := json.NewDecoder(dueResp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	for _, p := range body.Items {
		if p.ID == created.ID {
			t.Fatalf("freshly done plan %s still listed as due", p.ID)
		}
	}
}

func TestPlanLifecycleHappyPath(t *testing.T) {
	fakeNow := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	clock := OpsClock{NowFunc: func() time.Time { return fakeNow }}
	store := newPlanStore()
	svc := &PlanService{store: store, clock: clock}
	ctx := context.Background()
	plan, err := svc.Create(ctx, Plan{WellID: "well-07", Subject: "年度大修", Owner: "zhao", IntervalDays: 365})
	if err != nil {
		t.Fatal(err)
	}
	due, err := svc.Due(ctx, fakeNow.Add(24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if len(due) != 0 {
		t.Fatalf("new plan due too early: %+v", due)
	}
	fakeNow = fakeNow.AddDate(1, 0, 0)
	due, err = svc.Due(ctx, fakeNow)
	if err != nil {
		t.Fatal(err)
	}
	if len(due) != 1 || due[0].ID != plan.ID {
		t.Fatalf("plan not due after interval: %+v", due)
	}
	if _, err := svc.MarkDone(ctx, plan.ID, plan.Revision); err != nil {
		t.Fatal(err)
	}
	due, err = svc.Due(ctx, fakeNow)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range due {
		if p.ID == plan.ID {
			t.Fatalf("completed plan still due: %+v", p)
		}
	}
}
