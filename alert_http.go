package main

import (
	"net/http"
	"strings"
)

func alertsListHandler(svc *AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusParam := strings.TrimSpace(r.URL.Query().Get("status"))
		status := AlertStatus(statusParam)
		if statusParam == "" {
			status = AlertOpen
		}
		items := svc.List(status)
		if len(items) > 200 {
			items = items[:200]
		}
		opsJSON(w, http.StatusOK, map[string]any{"count": len(items), "open": svc.OpenCount(), "items": items})
	}
}

func alertAckHandler(svc *AlertService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		actor := opsActorFromRequest(r)
		alert, err := svc.Ack(r.Context(), id, actor)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, alert)
	}
}
