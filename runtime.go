package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

var requestSequence uint64

type serverLifecycle struct {
	inflight sync.WaitGroup
	active   atomic.Int64
}

func (l *serverLifecycle) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l.inflight.Add(1)
		l.active.Add(1)
		defer func() {
			l.inflight.Done()
			l.active.Add(-1)
		}()
		next.ServeHTTP(w, r)
	})
}

func (l *serverLifecycle) Active() int64 { return l.active.Load() }

func serveAddress(address string, handler http.Handler) error {
	lifecycle := &serverLifecycle{}
	server := newEnterpriseServer(address, handler, lifecycle)
	return serveHTTP(server, lifecycle, LoadConfig().ShutdownTimeoutSeconds)
}

func serveHTTP(server *http.Server, lifecycle *serverLifecycle, shutdownSeconds int) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signals)
	return serveHTTPListener(server, lifecycle, shutdownSeconds, nil, signals)
}

func serveHTTPListener(server *http.Server, lifecycle *serverLifecycle, shutdownSeconds int, ln net.Listener, signals <-chan os.Signal) error {
	errCh := make(chan error, 1)
	go func() {
		if ln != nil {
			errCh <- server.Serve(ln)
		} else {
			errCh <- server.ListenAndServe()
		}
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-signals:
		timeout := time.Duration(shutdownSeconds) * time.Second
		if timeout <= 0 {
			timeout = 10 * time.Second
		}
		shutdownContext, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		_ = server.Shutdown(shutdownContext)
		lifecycle.inflight.Wait()
		return nil
	}
}

func newEnterpriseServer(address string, handler http.Handler, lifecycle *serverLifecycle) *http.Server {
	return &http.Server{
		Addr: address,
		Handler: opsEnterpriseMiddleware(requestIDMiddleware(recoveryMiddleware(
			lifecycle.middleware(handler),
		))),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20,
	}
}

func requestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("req-%d", atomic.AddUint64(&requestSequence, 1))
		}
		w.Header().Set("X-Request-ID", requestID)
		next.ServeHTTP(w, r)
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic recovered request_id=%s method=%s path=%s panic=%v", w.Header().Get("X-Request-ID"), r.Method, r.URL.Path, recovered)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
