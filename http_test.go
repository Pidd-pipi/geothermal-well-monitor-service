package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPWorkflow(t *testing.T) {
	ts := httptest.NewServer(testRouter())
	defer ts.Close()
	for _, path := range []string{"/health", "/api/wells", "/"} {
		resp, err := http.Get(ts.URL + path)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("GET %s status=%v err=%v", path, resp.StatusCode, err)
		}
		resp.Body.Close()
	}
	resp, err := http.Post(ts.URL+"/api/wells/status", "application/json", bytes.NewBufferString(`{"id":"well-07","status":"isolated"}`))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("valid POST status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
	resp, err = http.Post(ts.URL+"/api/wells/status", "application/json", bytes.NewBufferString(`{"id":"well-07","status":"unknown"}`))
	if err != nil || resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
	resp, err = http.Post(ts.URL+"/api/wells/status", "application/json", bytes.NewBufferString(`{"id":"missing","status":"isolated"}`))
	if err != nil || resp.StatusCode == http.StatusOK {
		t.Fatalf("missing well not rejected: status=%v err=%v", resp.StatusCode, err)
	}
	resp.Body.Close()
}
