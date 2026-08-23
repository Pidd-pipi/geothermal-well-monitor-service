package main

import (
	"errors"
	"sync"
)

var ErrWellNotFound = errors.New("well not found")

type WellStore struct {
	mu    sync.RWMutex
	wells map[string]Well
}

func NewWellStore() *WellStore {
	return &WellStore{wells: map[string]Well{
		"well-07": {ID: "well-07", Name: "北场注采井 07", DepthM: 1840, PressureKPa: 812.4, TemperatureC: 142.6, Status: "producing"},
		"well-12": {ID: "well-12", Name: "东场生产井 12", DepthM: 2210, PressureKPa: 768.1, TemperatureC: 136.9, Status: "inspection"},
	}}
}
func (s *WellStore) List() []Well {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Well, 0, len(s.wells))
	for _, w := range s.wells {
		out = append(out, w)
	}
	return out
}
func (s *WellStore) UpdateStatus(id, status string) (Well, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	w, ok := s.wells[id]
	if !ok {
		return Well{}, ErrWellNotFound
	}
	w.Status = status
	s.wells[id] = w
	return w, nil
}
