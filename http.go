package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type statusRequest struct {
	ID     string `json:"id"`
	Status string `json:"status"`
}

type AppDeps struct {
	wells     *WellService
	readings  *ReadingService
	alerts    *AlertService
	tasks     *TaskService
	plans     *PlanService
	ops       *OpsService
	backfill  *BackfillRunner
	operators *OperatorStore
	sessions  *SessionStore
}

func NewRouter(deps *AppDeps) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", healthHandler)

	mux.HandleFunc("GET /api/wells", func(w http.ResponseWriter, r *http.Request) { writeJSON(w, http.StatusOK, deps.wells.Wells()) })
	mux.HandleFunc("POST /api/wells/status", func(w http.ResponseWriter, r *http.Request) {
		var req statusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || strings.TrimSpace(req.ID) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id and status JSON are required"})
			return
		}
		well, err := deps.wells.ChangeStatus(req.ID, req.Status)
		if errors.Is(err, ErrWellNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, well)
	})

	mux.HandleFunc("POST /api/readings", readingIngestHandler(deps.readings))
	mux.HandleFunc("GET /api/readings/recent", readingRecentHandler(deps.readings))
	mux.HandleFunc("GET /api/readings/stats", readingStatsHandler(deps.readings))

	mux.HandleFunc("GET /api/alerts", alertsListHandler(deps.alerts))
	mux.HandleFunc("POST /api/alerts/{id}/ack", alertAckHandler(deps.alerts))

	mux.HandleFunc("GET /api/tasks", tasksListHandler(deps.tasks))
	mux.HandleFunc("POST /api/tasks", taskCreateHandler(deps.tasks))
	mux.HandleFunc("POST /api/tasks/{id}/assign", taskAssignHandler(deps.tasks))
	mux.HandleFunc("POST /api/tasks/{id}/start", taskStartHandler(deps.tasks))
	mux.HandleFunc("POST /api/tasks/{id}/complete", taskCompleteHandler(deps.tasks))

	mux.HandleFunc("GET /api/plans", plansListHandler(deps.plans))
	mux.HandleFunc("POST /api/plans", planCreateHandler(deps.plans))
	mux.HandleFunc("POST /api/plans/{id}/done", planDoneHandler(deps.plans))

	mux.HandleFunc("GET /api/ops/records", opsRecordsListHandler(deps.ops))
	mux.HandleFunc("POST /api/ops/records", opsRecordCreateHandler(deps.ops))
	mux.HandleFunc("POST /api/ops/records/{id}/transition", opsRecordTransitionHandler(deps.ops))
	mux.HandleFunc("GET /api/ops/records/{id}/audit", opsRecordAuditHandler(deps.ops))
	mux.HandleFunc("GET /api/ops/snapshot", opsSnapshotHandler(deps.ops))

	mux.HandleFunc("POST /api/backfill", backfillStartHandler(deps.backfill))
	mux.HandleFunc("GET /api/backfill/{id}", backfillStatusHandler(deps.backfill))

	mux.HandleFunc("POST /api/operators/login", loginHandler(deps.operators, deps.sessions))
	mux.Handle("GET /api/operators/me", requireAuth(deps.sessions)(meHandler()))

	return withStatic(mux)
}
