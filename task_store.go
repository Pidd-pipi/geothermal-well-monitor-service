package main

import (
	"context"
	"sort"
	"sync"
)

type TaskStore struct {
	mu         sync.RWMutex
	items      map[string]Task
	assigning  map[string]bool
	completed  []Task
	maxHistory int
}

func newTaskStore() *TaskStore {
	return &TaskStore{items: map[string]Task{}, assigning: map[string]bool{}, completed: []Task{}, maxHistory: 100}
}

func (s *TaskStore) Put(ctx context.Context, t Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[t.ID]; ok {
		return ErrOpsConflict
	}
	s.items[t.ID] = normalizeTask(t)
	return nil
}

func (s *TaskStore) Get(ctx context.Context, id string) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.items[id]
	if !ok {
		return Task{}, ErrOpsNotFound
	}
	return t.Clone(), nil
}

func (s *TaskStore) List(ctx context.Context) ([]Task, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.items))
	for _, t := range s.items {
		out = append(out, t.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *TaskStore) Update(ctx context.Context, t Task, expected int) (Task, error) {
	select {
	case <-ctx.Done():
		return Task{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[t.ID]
	if !ok {
		return Task{}, ErrOpsNotFound
	}
	if expected > 0 && current.Revision != expected {
		return Task{}, ErrOpsConflict
	}
	t.Revision = current.Revision + 1
	s.items[t.ID] = t.Clone()
	return t.Clone(), nil
}

func (s *TaskStore) BeginAssign(ctx context.Context, id string) (bool, error) {
	select {
	case <-ctx.Done():
		return false, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return false, ErrOpsNotFound
	}
	if s.assigning[id] {
		return false, ErrOpsConflict
	}
	s.assigning[id] = true
	return true, nil
}

func (s *TaskStore) EndAssign(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.assigning, id)
	return nil
}

func (s *TaskStore) RecordCompletion(ctx context.Context, t Task) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.completed = append(s.completed, t.Clone())
	return nil
}

func (s *TaskStore) CompletedHistory() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Task, 0, len(s.completed))
	for _, t := range s.completed {
		out = append(out, t)
	}
	return out
}

func (s *TaskStore) AssigningCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.assigning)
}
