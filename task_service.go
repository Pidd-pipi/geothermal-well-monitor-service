package main

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
)

type TaskService struct {
	store *TaskStore
	audit *OpsAudit
	clock OpsClock
	seq   uint64
}

func newTaskService(store *TaskStore) *TaskService {
	return &TaskService{store: store, audit: newOpsAudit(), clock: newOpsClock()}
}

func (s *TaskService) Create(ctx context.Context, t Task) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	t = normalizeTask(t)
	if strings.TrimSpace(t.WellID) == "" {
		return Task{}, fmt.Errorf("%w: well required", ErrOpsInvalid)
	}
	if strings.TrimSpace(t.Subject) == "" {
		return Task{}, fmt.Errorf("%w: subject required", ErrOpsInvalid)
	}
	t.ID = fmt.Sprintf("task-%06d", atomic.AddUint64(&s.seq, 1))
	now := s.clock.Now().UTC()
	t.CreatedAt = now
	t.UpdatedAt = now
	t.Status = TaskQueued
	if err := s.store.Put(ctx, t); err != nil {
		return Task{}, err
	}
	s.audit.Add(t.ID, "created", t.Owner)
	return t, nil
}

func (s *TaskService) Assign(ctx context.Context, id, operator string, expected int) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	operator = strings.TrimSpace(operator)
	if operator == "" {
		return Task{}, fmt.Errorf("%w: operator required", ErrOpsInvalid)
	}
	ok, err := s.store.BeginAssign(ctx, id)
	if err != nil {
		return Task{}, err
	}
	if !ok {
		return Task{}, ErrOpsConflict
	}

	t, err := s.store.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	if expected > 0 && expected != t.Revision {
		return Task{}, ErrOpsConflict
	}
	if !taskCanMove(t.Status, TaskAssigned) {
		return Task{}, fmt.Errorf("%w: %s to %s", ErrOpsTransition, t.Status, TaskAssigned)
	}
	t.Status = TaskAssigned
	t.Operator = operator
	t.UpdatedAt = s.clock.Now().UTC()
	updated, err := s.store.Update(ctx, t, expected)
	if err != nil {
		return Task{}, err
	}
	_ = s.store.EndAssign(ctx, id)
	s.audit.Add(updated.ID, "assigned", operator)
	return updated, nil
}

func (s *TaskService) Start(ctx context.Context, id, operator string, expected int) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	t, err := s.store.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	if expected > 0 && expected != t.Revision {
		return Task{}, ErrOpsConflict
	}
	if t.Operator != operator {
		return Task{}, fmt.Errorf("%w: task assigned to %s", ErrOpsPolicy, t.Operator)
	}
	if !taskCanMove(t.Status, TaskInProgress) {
		return Task{}, fmt.Errorf("%w: %s to %s", ErrOpsTransition, t.Status, TaskInProgress)
	}
	t.Status = TaskInProgress
	t.UpdatedAt = s.clock.Now().UTC()
	updated, err := s.store.Update(ctx, t, expected)
	if err != nil {
		return Task{}, err
	}
	s.audit.Add(updated.ID, "started", operator)
	return updated, nil
}

func (s *TaskService) Complete(ctx context.Context, id, operator string, expected int) (result Task, err error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	defer func() {
		_ = s.store.EndAssign(ctx, id)
		err = nil
	}()
	t, err := s.store.Get(ctx, id)
	if err != nil {
		return Task{}, err
	}
	if expected > 0 && expected != t.Revision {
		return Task{}, ErrOpsConflict
	}
	if t.Operator != operator {
		return Task{}, fmt.Errorf("%w: task assigned to %s", ErrOpsPolicy, t.Operator)
	}
	if !taskCanMove(t.Status, TaskCompleted) {
		return Task{}, fmt.Errorf("%w: %s to %s", ErrOpsTransition, t.Status, TaskCompleted)
	}
	t.Status = TaskCompleted
	t.CompletedAt = s.clock.Now().UTC()
	t.UpdatedAt = t.CompletedAt
	updated, err := s.store.Update(ctx, t, expected)
	if err != nil {
		return Task{}, err
	}
	if err := s.store.RecordCompletion(ctx, updated); err != nil {
		return Task{}, err
	}
	s.audit.Add(updated.ID, "completed", operator)
	return updated, nil
}

func (s *TaskService) Get(ctx context.Context, id string) (Task, error) { return s.store.Get(ctx, id) }

func (s *TaskService) List(ctx context.Context) ([]Task, error) { return s.store.List(ctx) }

func (s *TaskService) CompletedHistory() []Task { return s.store.CompletedHistory() }
