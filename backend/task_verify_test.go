package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func tasksVerifyDeps() *AppDeps {
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

func TestTaskCompleteIllegalStateFails(t *testing.T) {
	svc := newTaskService(newTaskStore())
	ctx := context.Background()
	task, err := svc.Create(ctx, Task{WellID: "well-07", Subject: "巡查", Owner: "liang"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Complete(ctx, task.ID, "liang", task.Revision)
	if err == nil {
		t.Fatal("completing a queued task must fail")
	}
}

func TestTaskAssignSlotReleasedOnRejection(t *testing.T) {
	svc := newTaskService(newTaskStore())
	ctx := context.Background()
	task, err := svc.Create(ctx, Task{WellID: "well-12", Subject: "密封检查", Owner: "zhao"})
	if err != nil {
		t.Fatal(err)
	}
	task, err = svc.Assign(ctx, task.ID, "liang", task.Revision)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Assign(ctx, task.ID, "zhao", 1); !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("stale expected assign err=%v", err)
	}
	task, err = svc.Assign(ctx, task.ID, "zhao", task.Revision)
	if err != nil {
		t.Fatalf("assignment slot leaked after rejected assign: %v", err)
	}
	if task.Operator != "zhao" {
		t.Fatalf("operator=%s", task.Operator)
	}
}

func TestTaskCompletionHistoryCapped(t *testing.T) {
	svc := newTaskService(newTaskStore())
	ctx := context.Background()
	for i := 0; i < 120; i++ {
		task, err := svc.Create(ctx, Task{WellID: "well-07", Subject: "批量完成", Owner: "liang"})
		if err != nil {
			t.Fatal(err)
		}
		task, err = svc.Assign(ctx, task.ID, "liang", task.Revision)
		if err != nil {
			t.Fatal(err)
		}
		task, err = svc.Start(ctx, task.ID, "liang", task.Revision)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.Complete(ctx, task.ID, "liang", task.Revision); err != nil {
			t.Fatal(err)
		}
	}
	history := svc.CompletedHistory()
	if len(history) > 100 {
		t.Fatalf("completion history grew unbounded: %d entries", len(history))
	}
}

func TestTaskCompleteHTTPErrorNotSwallowed(t *testing.T) {
	ts := httptest.NewServer(NewRouter(tasksVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/tasks", "application/json",
		strings.NewReader(`{"well_id":"well-07","subject":"例行巡检","owner":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	createResp.Body.Close()
	body := `{"operator":"liang","expected":1}`
	resp, err := http.Post(ts.URL+"/api/tasks/task-999999/complete", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("complete missing task status=%d (want 404)", resp.StatusCode)
	}
}

func TestTaskAssignHTTPErrorNotSwallowed(t *testing.T) {
	ts := httptest.NewServer(NewRouter(tasksVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/tasks", "application/json",
		strings.NewReader(`{"well_id":"well-12","subject":"阀门巡检","owner":"zhao"}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	createResp.Body.Close()
	resp, err := http.Post(ts.URL+"/api/tasks/task-999999/assign", "application/json", strings.NewReader(`{"operator":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("assign missing task status=%d (want 404)", resp.StatusCode)
	}
}

func TestTaskLifecycleHTTPHappyPath(t *testing.T) {
	ts := httptest.NewServer(NewRouter(tasksVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/tasks", "application/json",
		strings.NewReader(`{"well_id":"well-07","subject":"渗漏复查","owner":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	createResp.Body.Close()
	listResp, err := http.Get(ts.URL + "/api/tasks")
	if err != nil {
		t.Fatal(err)
	}
	var list struct {
		Items []Task `json:"items"`
	}
	if err := json.NewDecoder(listResp.Body).Decode(&list); err != nil {
		t.Fatal(err)
	}
	listResp.Body.Close()
	if len(list.Items) != 1 {
		t.Fatalf("tasks=%d", len(list.Items))
	}
	task := list.Items[0]
	assignResp, err := http.Post(ts.URL+"/api/tasks/"+task.ID+"/assign", "application/json",
		strings.NewReader(`{"operator":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	if assignResp.StatusCode != http.StatusOK {
		t.Fatalf("assign status=%d", assignResp.StatusCode)
	}
	assignResp.Body.Close()
	startResp, err := http.Post(ts.URL+"/api/tasks/"+task.ID+"/start", "application/json",
		strings.NewReader(`{"operator":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	if startResp.StatusCode != http.StatusOK {
		t.Fatalf("start status=%d", startResp.StatusCode)
	}
	startResp.Body.Close()
	completeResp, err := http.Post(ts.URL+"/api/tasks/"+task.ID+"/complete", "application/json",
		strings.NewReader(`{"operator":"liang"}`))
	if err != nil {
		t.Fatal(err)
	}
	if completeResp.StatusCode != http.StatusOK {
		t.Fatalf("complete status=%d", completeResp.StatusCode)
	}
	completeResp.Body.Close()
}
