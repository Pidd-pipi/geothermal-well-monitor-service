package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func opsEnterpriseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Set("X-Operations-Domain", opsDomainName)
		if strings.TrimSpace(r.Header.Get("X-Request-ID")) == "" {
			w.Header().Set("X-Operations-Request", "generated")
		} else {
			w.Header().Set("X-Operations-Request", "provided")
		}
		defer func() { w.Header().Set("X-Operations-Latency-Ms", formatOpsInt(int(time.Since(start).Milliseconds()))) }()
		next.ServeHTTP(w, r)
	})
}
func formatOpsInt(value int) string {
	if value == 0 {
		return "0"
	}
	out := ""
	for value > 0 {
		out = string(rune('0'+value%10)) + out
		value /= 10
	}
	return out
}
func opsJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
func opsAllowed(method string, allowed ...string) bool {
	for _, candidate := range allowed {
		if method == candidate {
			return true
		}
	}
	return false
}
func opsPathID(path, prefix string) string {
	if !strings.HasPrefix(path, prefix) {
		return ""
	}
	return strings.Trim(strings.TrimPrefix(path, prefix), "/")
}
func opsActorFromRequest(r *http.Request) string {
	value := strings.TrimSpace(r.Header.Get("X-Operator"))
	if value == "" {
		return "web"
	}
	return value
}
func opsNoStore(w http.ResponseWriter)    { w.Header().Set("Cache-Control", "no-store") }
func opsRequestID(r *http.Request) string { return r.Header.Get("X-Request-ID") }

func writeOpsError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrOpsNotFound):
		opsJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrOpsConflict):
		opsJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrOpsTransition):
		opsJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrOpsInvalid):
		opsJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
	case errors.Is(err, ErrOpsPolicy):
		opsJSON(w, http.StatusForbidden, map[string]string{"error": err.Error()})
	default:
		opsJSON(w, http.StatusInternalServerError, map[string]string{"error": "internal error"})
	}
}

func opsRecordsListHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
		q := OpsQuery{
			Subject:  r.URL.Query().Get("subject"),
			Status:   OpsStatus(r.URL.Query().Get("status")),
			Priority: OpsPriority(r.URL.Query().Get("priority")),
			Owner:    r.URL.Query().Get("owner"),
			Page:     page,
			PageSize: size,
		}
		result, err := svc.Search(r.Context(), q)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, result)
	}
}

func opsRecordCreateHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Subject  string            `json:"subject"`
			Owner    string            `json:"owner"`
			Priority OpsPriority       `json:"priority"`
			Labels   map[string]string `json:"labels"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid record JSON"})
			return
		}
		record, err := svc.Create(r.Context(), OpsRecord{
			Subject:  body.Subject,
			Owner:    body.Owner,
			Priority: body.Priority,
			Labels:   body.Labels,
		})
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusCreated, record)
	}
}

func opsRecordTransitionHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body struct {
			Expected int       `json:"expected"`
			Target   OpsStatus `json:"target"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid transition JSON"})
			return
		}
		if body.Target == "" || !opsStatusValid(body.Target) {
			writeOpsError(w, fmt.Errorf("%w: invalid target status", ErrOpsInvalid))
			return
		}
		actor := opsActorFromRequest(r)
		record, err := svc.Transition(r.Context(), id, body.Expected, body.Target, actor)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, record)
	}
}

func opsRecordAuditHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		events := svc.Audit(id)
		opsJSON(w, http.StatusOK, map[string]any{"record_id": id, "count": len(events), "events": events})
	}
}

func opsSnapshotHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		opsJSON(w, http.StatusOK, svc.Snapshot())
	}
}
