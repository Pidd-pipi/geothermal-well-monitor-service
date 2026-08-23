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

func TestOpsTransitionInvalidTarget400(t *testing.T) {
	ts := httptest.NewServer(NewRouter(opsVerifyDeps()))
	defer ts.Close()
	createResp, err := http.Post(ts.URL+"/api/ops/records", "application/json",
		strings.NewReader(`{"subject":"压力异常","owner":"liang","priority":"high","labels":{"site":"north"}}`))
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
	body := fmt.Sprintf(`{"target":"bogus","expected":%d}`, created.Revision)
	resp, err := http.Post(ts.URL+"/api/ops/records/"+created.ID+"/transition", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid target status=%d (want 400)", resp.StatusCode)
	}
}

func TestOpsTransitionMoveErrorWrapped(t *testing.T) {
	svc := newOpsService(nil)
	ctx := context.Background()
	rec, err := svc.Create(ctx, OpsRecord{Subject: "闭环演练", Owner: "zhao", Priority: OpsPriorityNormal, Labels: map[string]string{"site": "east"}})
	if err != nil {
		t.Fatal(err)
	}
	// queued -> closed -> active: second move is illegal
	rec, err = svc.Transition(ctx, rec.ID, rec.Revision, OpsStatusClosed, "zhao")
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.Transition(ctx, rec.ID, rec.Revision, OpsStatusActive, "zhao")
	if err == nil {
		t.Fatal("closed -> active should be rejected")
	}
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("move error not wrapped as OpsError: %v", err)
	}
	if typed.Code != "transition" {
		t.Fatalf("move error code=%s", typed.Code)
	}
}
