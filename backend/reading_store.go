package main

import (
	"sort"
	"sync"
	"time"
)

type ReadingStore struct {
	mu      sync.RWMutex
	byWell  map[string][]Reading
	maxKeep int
}

func newReadingStore(maxKeep int) *ReadingStore {
	if maxKeep < 1 {
		maxKeep = 1000
	}
	return &ReadingStore{byWell: map[string][]Reading{}, maxKeep: maxKeep}
}

func (s *ReadingStore) Append(r Reading) {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.byWell[r.WellID]
	items = append(items, r)
	if len(items) > s.maxKeep {
		items = append([]Reading(nil), items[len(items)-s.maxKeep:]...)
	}
	s.byWell[r.WellID] = items
}

func (s *ReadingStore) Recent(wellID string, n int) []Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.byWell[wellID]
	if n < 1 {
		n = 10
	}
	if n > len(items) {
		n = len(items)
	}
	out := make([]Reading, 0, n)
	for i := len(items) - n; i < len(items); i++ {
		out = append(out, items[i])
	}
	return out
}

func (s *ReadingStore) Range(wellID string, from, to time.Time) []Reading {
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := s.byWell[wellID]
	out := []Reading{}
	for _, item := range items {
		if !item.At.Before(from) && item.At.Before(to) {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].At.Before(out[j].At) })
	return out
}

func (s *ReadingStore) Prune(wellID string, keep int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	items := s.byWell[wellID]
	if keep < 1 {
		keep = s.maxKeep
	}
	if len(items) <= keep {
		return 0
	}
	removed := len(items) - keep
	s.byWell[wellID] = append([]Reading(nil), items[len(items)-keep:]...)
	return removed
}

func (s *ReadingStore) Count(wellID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.byWell[wellID])
}
