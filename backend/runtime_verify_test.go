package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"
)

func TestOpsContextParentCancel(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	ctx, _ := opsContext(parent, 30*time.Second)
	cancel()
	select {
	case <-ctx.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("opsContext did not propagate parent cancellation")
	}
}

func TestOpsDelayHonorsCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- opsDelay(ctx, 5*time.Second) }()
	time.Sleep(50 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("opsDelay err=%v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("opsDelay ignored cancellation")
	}
}

func TestServerShutdownWaitsForInflight(t *testing.T) {
	lifecycle := &serverLifecycle{}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/slow" {
			time.Sleep(4 * time.Second)
		}
		w.WriteHeader(http.StatusOK)
	})
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	server := &http.Server{Addr: ln.Addr().String(), Handler: lifecycle.middleware(handler)}
	signals := make(chan os.Signal, 1)
	ret := make(chan error, 1)
	go func() { ret <- serveHTTPListener(server, lifecycle, 1, ln, signals) }()
	ready := false
	for i := 0; i < 100; i++ {
		resp, err := http.Get("http://" + ln.Addr().String() + "/")
		if err == nil {
			resp.Body.Close()
			ready = true
			break
		}
		time.Sleep(20 * time.Millisecond)
	}
	if !ready {
		t.Fatal("server not ready")
	}
	slowDone := make(chan struct{})
	go func() {
		resp, err := http.Get("http://" + ln.Addr().String() + "/slow")
		if err == nil {
			resp.Body.Close()
		}
		close(slowDone)
	}()
	inflight := false
	for i := 0; i < 200; i++ {
		if lifecycle.Active() == 1 {
			inflight = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !inflight {
		t.Fatal("/slow request never reached handler")
	}
	signals <- syscall.SIGTERM
	select {
	case err := <-ret:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(8 * time.Second):
		t.Fatal("serveHTTPListener did not return")
	}
	if lifecycle.Active() != 0 {
		t.Fatalf("serveHTTPListener returned while %d in-flight handlers still running", lifecycle.Active())
	}
	<-slowDone
}

func TestLoadConfigShutdownTimeoutHonorsEnv(t *testing.T) {
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "5")
	cfg := LoadConfig()
	if cfg.ShutdownTimeoutSeconds != 5 {
		t.Fatalf("shutdown timeout=%d (want 5)", cfg.ShutdownTimeoutSeconds)
	}
}

func TestLoadConfigRequestTimeoutFromEnv(t *testing.T) {
	t.Setenv("REQUEST_TIMEOUT_SECONDS", "2")
	cfg := LoadConfig()
	if cfg.RequestTimeoutSeconds != 2 {
		t.Fatalf("request timeout=%d (want 2)", cfg.RequestTimeoutSeconds)
	}
}
