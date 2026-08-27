package main

import "time"

type AlertKind string

const (
	AlertPressureHigh AlertKind = "pressure_high"
	AlertPressureLow  AlertKind = "pressure_low"
	AlertTempHigh     AlertKind = "temperature_high"
)

type AlertStatus string

const (
	AlertOpen  AlertStatus = "open"
	AlertAcked AlertStatus = "acknowledged"
	AlertSeen  AlertStatus = "seen"
)

type Alert struct {
	ID        string
	WellID    string
	Kind      AlertKind
	Severity  OpsPriority
	Threshold float64
	Value     float64
	At        time.Time
	Status    AlertStatus
	Actor     string
}

func (a Alert) Clone() Alert { return a }

type AlertRule struct {
	Kind      AlertKind
	Severity  OpsPriority
	Threshold float64
	Above     bool
}

func alertRules() []AlertRule {
	return []AlertRule{
		{Kind: AlertPressureHigh, Severity: OpsPriorityHigh, Threshold: 1500, Above: true},
		{Kind: AlertPressureLow, Severity: OpsPriorityHigh, Threshold: 300, Above: false},
		{Kind: AlertTempHigh, Severity: OpsPriorityCritical, Threshold: 250, Above: true},
	}
}

func alertRuleFor(kind AlertKind) (AlertRule, bool) {
	for _, rule := range alertRules() {
		if rule.Kind == kind {
			return rule, true
		}
	}
	return AlertRule{}, false
}

func alertsForReading(r Reading) []Alert {
	out := []Alert{}
	value := 0.0
	for _, rule := range alertRules() {
		switch rule.Kind {
		case AlertPressureHigh:
			value = r.PressureKPa
		case AlertPressureLow:
			value = r.PressureKPa
		case AlertTempHigh:
			value = r.TemperatureC
		}
		if (rule.Above && value > rule.Threshold) || (!rule.Above && value < rule.Threshold) {
			out = append(out, Alert{
				WellID:    r.WellID,
				Kind:      rule.Kind,
				Severity:  rule.Severity,
				Threshold: rule.Threshold,
				Value:     value,
			})
		}
	}
	return out
}
