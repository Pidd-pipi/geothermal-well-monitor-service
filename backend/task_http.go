package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func tasksListHandler(svc *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items, err := svc.List(r.Context())
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"count": len(items), "items": items})
	}
}

func taskCreateHandler(svc *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			WellID   string            `json:"well_id"`
			Subject  string            `json:"subject"`
			Owner    string            `json:"owner"`
			Priority OpsPriority       `json:"priority"`
			Labels   map[string]string `json:"labels"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid task JSON"})
			return
		}
		task, err := svc.Create(r.Context(), Task{
			WellID:   body.WellID,
			Subject:  body.Subject,
			Owner:    body.Owner,
			Priority: body.Priority,
			Labels:   body.Labels,
		})
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusCreated, task)
	}
}

func taskAssignHandler(svc *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body struct {
			Operator string `json:"operator"`
			Expected int    `json:"expected"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid assign JSON"})
			return
		}
		task, err := svc.Assign(r.Context(), id, body.Operator, body.Expected)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, task)
	}
}

func taskStartHandler(svc *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body struct {
			Operator string `json:"operator"`
			Expected int    `json:"expected"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid start JSON"})
			return
		}
		task, err := svc.Start(r.Context(), id, body.Operator, body.Expected)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, task)
	}
}

func taskCompleteHandler(svc *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		var body struct {
			Operator string `json:"operator"`
			Expected int    `json:"expected"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid complete JSON"})
			return
		}
		task, err := svc.Complete(r.Context(), id, body.Operator, body.Expected)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"task": task, "completed": strconv.FormatBool(true)})
	}
}
