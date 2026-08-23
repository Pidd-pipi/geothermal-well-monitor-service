package main

import (
	"strings"
	"time"
)

type PlanStatus string

const (
	PlanActive   PlanStatus = "active"
	PlanPaused   PlanStatus = "paused"
	PlanArchived PlanStatus = "archived"
)

type Plan struct {
	ID           string
	WellID       string
	Subject      string
	Owner        string
	Status       PlanStatus
	IntervalDays int
	LastDoneAt   time.Time
	NextDueAt    time.Time
	Revision     int
	CreatedAt    time.Time
}

func (p Plan) Clone() Plan { return p }

func normalizePlan(p Plan) Plan {
	p.WellID = strings.ToLower(strings.TrimSpace(p.WellID))
	p.Subject = strings.Join(strings.Fields(p.Subject), " ")
	p.Owner = strings.TrimSpace(p.Owner)
	if p.IntervalDays < 1 {
		p.IntervalDays = 30
	}
	if p.Status == "" {
		p.Status = PlanActive
	}
	if p.Revision < 1 {
		p.Revision = 1
	}
	return p
}

func planNextDue(last time.Time, intervalDays int) time.Time {
	return last.AddDate(0, 0, intervalDays)
}

func planIsActive(p Plan) bool { return p.Status == PlanActive }

var planTransitionTable = map[PlanStatus]map[PlanStatus]bool{
	PlanActive:   {PlanPaused: true, PlanArchived: true},
	PlanPaused:   {PlanActive: true, PlanArchived: true},
	PlanArchived: {},
}

func planCanMove(from, to PlanStatus) bool {
	return from == to || planTransitionTable[from][to]
}
