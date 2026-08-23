package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func authVerifyDeps() *AppDeps {
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

func TestLoginUnknownOperatorRejected401(t *testing.T) {
	ts := httptest.NewServer(NewRouter(authVerifyDeps()))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/operators/login", "application/json", strings.NewReader(`{"username":"nobody","password":"x"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("unknown login status=%d (want 401)", resp.StatusCode)
	}
}

func TestLoginIssuesSessionForValidUser(t *testing.T) {
	ts := httptest.NewServer(NewRouter(authVerifyDeps()))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/operators/login", "application/json", strings.NewReader(`{"username":"liang","password":"well2026"}`))
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("valid login status=%d", resp.StatusCode)
	}
	var out struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if out.Token == "" || out.Username != "liang" {
		t.Fatalf("login response %+v", out)
	}
	req, _ := http.NewRequest("GET", ts.URL+"/api/operators/me", nil)
	req.Header.Set("X-Session-Token", out.Token)
	meResp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer meResp.Body.Close()
	if meResp.StatusCode != http.StatusOK {
		t.Fatalf("me with valid token status=%d", meResp.StatusCode)
	}
	var me map[string]any
	if err := json.NewDecoder(meResp.Body).Decode(&me); err != nil {
		t.Fatal(err)
	}
	if me["username"] != "liang" {
		t.Fatalf("me response %v", me)
	}
}

func TestMeEndpointRejectsInvalidToken(t *testing.T) {
	ts := httptest.NewServer(NewRouter(authVerifyDeps()))
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/api/operators/me", nil)
	req.Header.Set("X-Session-Token", "bogus-token")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("invalid token status=%d (want 401)", resp.StatusCode)
	}
}

func TestMeEndpointRejectsSpoofedOperator(t *testing.T) {
	ts := httptest.NewServer(NewRouter(authVerifyDeps()))
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/api/operators/me", nil)
	req.Header.Set("X-Operator", "admin")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("spoofed operator status=%d (want 401)", resp.StatusCode)
	}
}
