package main

import (
	"sort"
	"strings"
)

const opsDomainName = "geothermal-well-monitor-service"

type OpsStatus string

const (
	OpsStatusQueued OpsStatus = "queued"
	OpsStatusActive OpsStatus = "active"
	OpsStatusPaused OpsStatus = "paused"
	OpsStatusClosed OpsStatus = "closed"
)

type OpsPriority string

const (
	OpsPriorityLow      OpsPriority = "low"
	OpsPriorityNormal   OpsPriority = "normal"
	OpsPriorityHigh     OpsPriority = "high"
	OpsPriorityCritical OpsPriority = "critical"
)

type OpsRecord struct {
	ID        string
	Subject   string
	Owner     string
	Status    OpsStatus
	Priority  OpsPriority
	Revision  int
	Labels    map[string]string
	CreatedAt string
	UpdatedAt string
}

type OpsRule struct {
	Code           string
	Name           string
	Severity       OpsPriority
	RequiredLabels []string
	Terminal       bool
}

type OpsEvent struct {
	ID       string
	RecordID string
	Type     string
	Actor    string
	At       string
	Details  map[string]string
}

type OpsQuery struct {
	Subject  string
	Status   OpsStatus
	Priority OpsPriority
	Owner    string
	Page     int
	PageSize int
}

type OpsPage struct {
	Items    []OpsRecord
	Page     int
	PageSize int
	Total    int
	HasNext  bool
}

type OpsSnapshot struct {
	Domain      string
	GeneratedAt string
	Records     int
	Active      int
	ByStatus    map[OpsStatus]int
	ByPriority  map[OpsPriority]int
}

func (r OpsRecord) Clone() OpsRecord {
	copy := r
	copy.Labels = r.Labels
	return copy
}

func (r OpsRecord) LabelValue(key string) string { return r.Labels[key] }
func (r OpsRecord) Terminal() bool               { return r.Status == OpsStatusClosed }

func (p OpsPriority) Weight() int {
	switch p {
	case OpsPriorityCritical:
		return 4
	case OpsPriorityHigh:
		return 3
	case OpsPriorityNormal:
		return 2
	default:
		return 1
	}
}

func normalizeOpsRecord(record OpsRecord) OpsRecord {
	record.ID = strings.ToLower(strings.TrimSpace(record.ID))
	record.Subject = strings.Join(strings.Fields(record.Subject), " ")
	record.Owner = strings.TrimSpace(record.Owner)
	if record.Revision < 1 {
		record.Revision = 1
	}
	if record.Status == "" {
		record.Status = OpsStatusQueued
	}
	if record.Labels == nil {
		record.Labels = map[string]string{}
	}
	return record
}

func sortOpsRecords(items []OpsRecord) {
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Priority.Weight() != items[j].Priority.Weight() {
			return items[i].Priority.Weight() > items[j].Priority.Weight()
		}
		return items[i].UpdatedAt > items[j].UpdatedAt
	})
}

func opsRules() []OpsRule {
	out := make([]OpsRule, 0, 112)
	for _, group := range [][]OpsRule{
		opsRules01(), opsRules02(), opsRules03(), opsRules04(), opsRules05(), opsRules06(), opsRules07(),
		opsRules08(), opsRules09(), opsRules10(), opsRules11(), opsRules12(), opsRules13(), opsRules14(),
	} {
		out = append(out, group...)
	}
	return out
}
