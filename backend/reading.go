package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	pressureMinKPa = 100.0
	pressureMaxKPa = 2000.0
	tempMinC       = 20.0
	tempMaxC       = 300.0
)

type Reading struct {
	ID           string
	WellID       string
	At           time.Time
	PressureKPa  float64
	TemperatureC float64
	Source       string
}

func (r Reading) Clone() Reading { return r }

func normalizeReading(r Reading) (Reading, error) {
	r.WellID = strings.ToLower(strings.TrimSpace(r.WellID))
	r.Source = strings.TrimSpace(r.Source)
	if r.Source == "" {
		r.Source = "sensor"
	}
	if r.WellID == "" {
		return Reading{}, fmt.Errorf("%w: well_id required", ErrOpsInvalid)
	}
	if r.PressureKPa < pressureMinKPa || r.PressureKPa > pressureMaxKPa {
		return Reading{}, fmt.Errorf("%w: pressure %.1f out of range [%.1f, %.1f]", ErrOpsInvalid, r.PressureKPa, pressureMinKPa, pressureMaxKPa)
	}
	if r.TemperatureC < tempMinC || r.TemperatureC > tempMaxC {
		return Reading{}, fmt.Errorf("%w: temperature %.1f out of range [%.1f, %.1f]", ErrOpsInvalid, r.TemperatureC, tempMinC, tempMaxC)
	}
	if r.At.IsZero() {
		r.At = time.Now().UTC()
	}
	return r, nil
}

func readingSummary(r Reading) string {
	return fmt.Sprintf("%s %.1f %.1f", r.WellID, r.PressureKPa, r.TemperatureC)
}

func readingAge(now time.Time, r Reading) time.Duration {
	if now.Before(r.At) {
		return 0
	}
	return now.Sub(r.At)
}
