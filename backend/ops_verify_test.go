package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func opsVerifyDeps() *AppDeps {
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

func TestOpsTransitionMissingRecordHTTP404(t *testing.T) {
	ts := httptest.NewServer(NewRouter(opsVerifyDeps()))
	defer ts.Close()
	body := `{"target":"active","expected":1}`
	resp, err := http.Post(ts.URL+"/api/ops/records/no-such-record/transition", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("missing record transition status=%d (want 404)", resp.StatusCode)
	}
}

func TestOpsTransitionRevisionConflictHTTP409(t *testing.T) {
	ts := httptest.NewServer(NewRouter(opsVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/ops/records", "application/json",
		strings.NewReader(`{"subject":"井口渗漏检查","owner":"liang","priority":"high","labels":{"site":"north"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	var created OpsRecord
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	createResp.Body.Close()
	body := fmt.Sprintf(`{"target":"active","expected":%d}`, created.Revision+7)
	resp, err := http.Post(ts.URL+"/api/ops/records/"+created.ID+"/transition", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusConflict {
		t.Fatalf("revision conflict status=%d (want 409)", resp.StatusCode)
	}
}

func TestOpsCreateDuplicateWrappedConflict(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	rec := OpsRecord{ID: "rec-dup", Subject: "重复记录", Owner: "liang", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "north"}}
	if _, err := svc.Create(ctx, rec); err != nil {
		t.Fatal(err)
	}
	_, err := svc.Create(ctx, rec)
	if err == nil {
		t.Fatal("duplicate create accepted")
	}
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("duplicate create error not chainable to ErrOpsConflict: %v", err)
	}
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("duplicate create error not an OpsError: %v", err)
	}
}

func TestOpsWrappedErrorKeepsChain(t *testing.T) {
	err := wrapOps("transition", "store.get", ErrOpsNotFound)
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("wrapOps lost ErrOpsNotFound: %v", err)
	}
	var typed *OpsError
	if !errors.As(err, &typed) || typed.Code != "transition" {
		t.Fatalf("wrapOps result not unwrap-able OpsError: %v", err)
	}
}

func TestOpsTransitionHappyPathHTTP(t *testing.T) {
	ts := httptest.NewServer(NewRouter(opsVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/ops/records", "application/json",
		strings.NewReader(`{"subject":"阀门检修","owner":"zhao","priority":"normal","labels":{"site":"east"}}`))
	if err != nil {
		t.Fatal(err)
	}
	if createResp.StatusCode != http.StatusCreated {
		t.Fatalf("create status=%d", createResp.StatusCode)
	}
	var created OpsRecord
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatal(err)
	}
	createResp.Body.Close()
	body := fmt.Sprintf(`{"target":"active","expected":%d}`, created.Revision)
	resp, err := http.Post(ts.URL+"/api/ops/records/"+created.ID+"/transition", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("happy transition status=%d", resp.StatusCode)
	}
	var updated OpsRecord
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatal(err)
	}
	if updated.Status != OpsStatusActive || updated.Revision != created.Revision+1 {
		t.Fatalf("updated=%+v", updated)
	}
}
