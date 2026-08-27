package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func readingIngestHandler(svc *ReadingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			WellID       string  `json:"well_id"`
			PressureKPa  float64 `json:"pressure_kpa"`
			TemperatureC float64 `json:"temperature_c"`
			Source       string  `json:"source"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			opsJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid readings JSON"})
			return
		}
		reading, err := svc.Ingest(r.Context(), Reading{
			WellID:       body.WellID,
			PressureKPa:  body.PressureKPa,
			TemperatureC: body.TemperatureC,
			Source:       body.Source,
		})
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusCreated, map[string]any{"reading": reading, "summary": readingSummary(reading)})
	}
}

func readingRecentHandler(svc *ReadingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wellID := strings.TrimSpace(r.URL.Query().Get("well_id"))
		n, _ := strconv.Atoi(r.URL.Query().Get("n"))
		items, err := svc.Recent(r.Context(), wellID, n)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, map[string]any{"well_id": wellID, "count": len(items), "items": items})
	}
}

func readingStatsHandler(svc *ReadingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wellID := strings.TrimSpace(r.URL.Query().Get("well_id"))
		hours, _ := strconv.Atoi(r.URL.Query().Get("hours"))
		since := time.Time{}
		if hours > 0 {
			since = time.Now().UTC().Add(-time.Duration(hours) * time.Hour)
		}
		stats, err := svc.Stats(r.Context(), wellID, since)
		if err != nil {
			writeOpsError(w, err)
			return
		}
		opsJSON(w, http.StatusOK, stats)
	}
}
