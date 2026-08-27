package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

func plansListHandler(svc *PlanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("due") == "true" {
			items, err := svc.Due(r.Context(), time.Now().UTC())
			if err != nil {
				writeOpsError(w, err)
				return
			}
			opsJSON(w, http.StatusOK, map[string]any{"due": true, "count": len(items), "items": items})
			return
		}
		items, err := svc.List(r.Context())
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"count": len(items), "items": items})
	}
}

func planCreateHandler(svc *PlanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			WellID       string `json:"well_id"`
			Subject      string `json:"subject"`
			Owner        string `json:"owner"`
			IntervalDays int    `json:"interval_days"`
			Status       string `json:"status"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid plan JSON"})
			return
		}
		plan, err := svc.Create(r.Context(), Plan{
			WellID:       body.WellID,
			Subject:      body.Subject,
			Owner:        body.Owner,
			IntervalDays: body.IntervalDays,
			Status:       PlanStatus(body.Status),
		})
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusCreated, plan)
	}
}

func planDoneHandler(svc *PlanService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body struct {
			Expected int `json:"expected"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body)
		plan, err := svc.MarkDone(r.Context(), id, body.Expected)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"plan": plan, "done": strconv.FormatBool(true)})
	}
}
