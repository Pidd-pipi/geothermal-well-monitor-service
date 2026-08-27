package main

import (
	"encoding/json"
	"net/http"
)

func backfillStartHandler(runner *BackfillRunner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Payload string `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid backfill JSON"})
			return
		}
		job, err := runner.Start(r.Context(), body.Payload)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusAccepted, map[string]any{"id": job.ID, "status": job.Status, "total": job.Total})
	}
}

func backfillStatusHandler(runner *BackfillRunner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		job, _ := runner.Get(r.Context(), id)
		opsJSON(w, http.StatusOK, job)
	}
}
