package main

import (
	"context"
	"fmt"
	"sync/atomic"
)

type AlertService struct {
	store *AlertStore
	clock OpsClock
	seq   uint64
}

func newAlertService(store *AlertStore) *AlertService {
	return &AlertService{store: store, clock: newOpsClock()}
}

func (s *AlertService) EvaluateReading(ctx context.Context, r Reading) ([]Alert, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	specs := alertsForReading(r)
	created := []Alert{}
	for _, spec := range specs {
		if s.store.HasOpen(spec.WellID, spec.Kind) {
			continue
		}
		spec.ID = fmt.Sprintf("alt-%06d", atomic.AddUint64(&s.seq, 1))
		spec.At = s.clock.Now().UTC()
		spec.Status = AlertOpen
		s.store.Add(spec)
		created = append(created, spec)
	}
	return created, nil
}

func (s *AlertService) Ack(ctx context.Context, id, actor string) (Alert, error) {
	select {
	case <-ctx.Done():
		return Alert{}, ctx.Err()
	default:
	}
	current, ok := s.store.Get(id)
	if !ok {
		return Alert{}, ErrOpsNotFound
	}
	if current.Status != AlertOpen {
		return Alert{}, fmt.Errorf("%w: alert %s is not open", ErrOpsConflict, id)
	}
	current.Status = AlertAcked
	current.Actor = actor
	if !s.store.Replace(id, current) {
		return Alert{}, ErrOpsNotFound
	}
	return current, nil
}

func (s *AlertService) List(status AlertStatus) []Alert { return s.store.List(status) }

func (s *AlertService) OpenCount() int {
	items := s.store.List(AlertOpen)
	return len(items)
}
