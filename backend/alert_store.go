package main

import "sync"

type AlertStore struct {
	mu     sync.RWMutex
	alerts []Alert
}

func newAlertStore() *AlertStore { return &AlertStore{alerts: []Alert{}} }

func (s *AlertStore) Add(a Alert) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts = append(s.alerts, a)
}

func (s *AlertStore) Get(id string) (Alert, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.alerts {
		if a.ID == id {
			return a, true
		}
	}
	return Alert{}, false
}

func (s *AlertStore) List(status AlertStatus) []Alert {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Alert, 0, len(s.alerts))
	for _, a := range s.alerts {
		if status == "" || a.Status == status {
			out = append(out, a)
		}
	}
	return out
}

func (s *AlertStore) HasOpen(wellID string, kind AlertKind) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, a := range s.alerts {
		if a.WellID == wellID && a.Kind == kind && a.Status == AlertOpen {
			return true
		}
	}
	return false
}

func (s *AlertStore) Replace(id string, a Alert) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.alerts {
		if s.alerts[i].ID == id {
			s.alerts[i] = a
			return true
		}
	}
	return false
}

func (s *AlertStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.alerts)
}
