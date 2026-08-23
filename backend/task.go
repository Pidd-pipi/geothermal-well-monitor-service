package main

import (
	"strings"
	"time"
)

type TaskStatus string

const (
	TaskQueued     TaskStatus = "queued"
	TaskAssigned   TaskStatus = "assigned"
	TaskInProgress TaskStatus = "in_progress"
	TaskCompleted  TaskStatus = "completed"
	TaskCancelled  TaskStatus = "cancelled"
)

type Task struct {
	ID          string
	WellID      string
	Subject     string
	Owner       string
	Operator    string
	Status      TaskStatus
	Priority    OpsPriority
	Revision    int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
	Labels      map[string]string
}

func (t Task) Clone() Task {
	c := t
	c.Labels = map[string]string{}
	for k, v := range t.Labels {
		c.Labels[k] = v
	}
	return c
}

func normalizeTask(t Task) Task {
	t.WellID = strings.ToLower(strings.TrimSpace(t.WellID))
	t.Subject = strings.Join(strings.Fields(t.Subject), " ")
	t.Operator = strings.TrimSpace(t.Operator)
	if t.Priority == "" {
		t.Priority = OpsPriorityNormal
	}
	if t.Revision < 1 {
		t.Revision = 1
	}
	if t.Labels == nil {
		t.Labels = map[string]string{}
	}
	return t
}

var taskTransitionTable = map[TaskStatus]map[TaskStatus]bool{
	TaskQueued:     {TaskAssigned: true, TaskCancelled: true},
	TaskAssigned:   {TaskInProgress: true, TaskCancelled: true},
	TaskInProgress: {TaskCompleted: true, TaskCancelled: true},
	TaskCompleted:  {},
	TaskCancelled:  {},
}

func taskCanMove(from, to TaskStatus) bool {
	return from == to || taskTransitionTable[from][to]
}

func taskStatusValid(value TaskStatus) bool {
	return value == TaskQueued || value == TaskAssigned || value == TaskInProgress || value == TaskCompleted || value == TaskCancelled
}
