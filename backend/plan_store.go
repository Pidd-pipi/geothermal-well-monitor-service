package main

import (
	"context"
	"sort"
	"sync"
)

type PlanStore struct {
	mu    sync.RWMutex
	items map[string]Plan
}

func newPlanStore() *PlanStore { return &PlanStore{items: map[string]Plan{}} }

func (s *PlanStore) Put(ctx context.Context, p Plan) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[p.ID]; ok {
		return ErrOpsConflict
	}
	s.items[p.ID] = normalizePlan(p)
	return nil
}

func (s *PlanStore) Get(ctx context.Context, id string) (Plan, error) {
	select {
	case <-ctx.Done():
		return Plan{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.items[id]
	if !ok {
		return Plan{}, ErrOpsNotFound
	}
	return p.Clone(), nil
}

func (s *PlanStore) List(ctx context.Context) ([]Plan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Plan, 0, len(s.items))
	for _, p := range s.items {
		out = append(out, p.Clone())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (s *PlanStore) Update(ctx context.Context, p Plan, expected int) (Plan, error) {
	select {
	case <-ctx.Done():
		return Plan{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[p.ID]
	if !ok {
		return Plan{}, ErrOpsNotFound
	}
	if expected > 0 && current.Revision != expected {
		return Plan{}, ErrOpsConflict
	}
	p.Revision = current.Revision + 1
	s.items[p.ID] = p.Clone()
	return p.Clone(), nil
}

func (s *PlanStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.items)
}
