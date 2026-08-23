package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func wellsVerifyDeps() *AppDeps {
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

type wellPayload struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

func postStatus(t *testing.T, ts *httptest.Server, payload string) *http.Response {
	t.Helper()
	resp, err := http.Post(ts.URL+"/api/wells/status", "application/json", bytes.NewBufferString(payload))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func listWells(t *testing.T, ts *httptest.Server) []Well {
	t.Helper()
	resp, err := http.Get(ts.URL + "/api/wells")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var wells []Well
	if err := json.NewDecoder(resp.Body).Decode(&wells); err != nil {
		t.Fatal(err)
	}
	return wells
}

func wellByID(wells []Well, id string) (Well, bool) {
	for _, w := range wells {
		if w.ID == id {
			return w, true
		}
	}
	return Well{}, false
}

func TestWellStatusInvalidNotPersisted(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"well-07","status":"bogus"}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid status code=%d (want 400)", resp.StatusCode)
	}
	resp.Body.Close()
	wells := listWells(t, ts)
	w, ok := wellByID(wells, "well-07")
	if !ok {
		t.Fatal("well-07 missing")
	}
	if w.Status == "bogus" {
		t.Fatalf("invalid status persisted into store: %+v", w)
	}
}

func TestWellStatusResponseFresh(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"well-07","status":"isolated"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid update code=%d", resp.StatusCode)
	}
	var updated Well
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if updated.Status != "isolated" {
		t.Fatalf("response shows stale status %q", updated.Status)
	}
	wells := listWells(t, ts)
	w, _ := wellByID(wells, "well-07")
	if w.Status != "isolated" {
		t.Fatalf("list status=%q", w.Status)
	}
}

func TestWellStatusEmptyRejected(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"well-07","status":""}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("empty status code=%d (want 400)", resp.StatusCode)
	}
	resp.Body.Close()
	wells := listWells(t, ts)
	w, _ := wellByID(wells, "well-07")
	if w.Status == "" {
		t.Fatalf("empty status persisted: %+v", w)
	}
}

func TestWellStatusInvalidOnMissingWell400(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"no-such-well","status":"bogus"}`)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid status on missing well code=%d (want 400)", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestWellStatusValidUpdate(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"well-12","status":"producing"}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid update code=%d", resp.StatusCode)
	}
	resp.Body.Close()
}

func TestWellStatusMissingWell404(t *testing.T) {
	ts := httptest.NewServer(NewRouter(wellsVerifyDeps()))
	defer ts.Close()
	resp := postStatus(t, ts, `{"id":"no-such-well","status":"isolated"}`)
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("missing well status code=%d (want 404)", resp.StatusCode)
	}
	resp.Body.Close()
}
