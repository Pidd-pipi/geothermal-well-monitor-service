package main

import (
	"context"
	"sort"
	"sync"
)

type OpsStore struct {
	mu    sync.RWMutex
	items map[string]OpsRecord
}

func newOpsStore(seed []OpsRecord) *OpsStore {
	s := &OpsStore{items: map[string]OpsRecord{}}
	for _, item := range seed {
		item = normalizeOpsRecord(item)
		s.items[item.ID] = item
	}
	return s
}
func (s *OpsStore) Get(ctx context.Context, id string) (OpsRecord, error) {
	select {
	case <-ctx.Done():
		return OpsRecord{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok {
		return OpsRecord{}, ErrOpsNotFound
	}
	return item, nil
}
func (s *OpsStore) List(ctx context.Context) ([]OpsRecord, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]OpsRecord, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
func (s *OpsStore) Put(ctx context.Context, item OpsRecord) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[item.ID]; ok {
		return ErrOpsConflict
	}
	s.items[item.ID] = item
	return nil
}
func (s *OpsStore) Update(ctx context.Context, item OpsRecord, expected int) (OpsRecord, error) {
	select {
	case <-ctx.Done():
		return OpsRecord{}, ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	current, ok := s.items[item.ID]
	if !ok {
		return OpsRecord{}, ErrOpsNotFound
	}
	if expected > 0 && current.Revision != expected {
		return OpsRecord{}, ErrOpsConflict
	}
	item.Revision = current.Revision + 1
	item.UpdatedAt = timeNowOps()
	s.items[item.ID] = item
	return item, nil
}
func (s *OpsStore) Delete(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[id]; !ok {
		return ErrOpsNotFound
	}
	delete(s.items, id)
	return nil
}
func (s *OpsStore) Count() int { s.mu.RLock(); defer s.mu.RUnlock(); return len(s.items) }
