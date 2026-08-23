package main

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

type PlanService struct {
	store *PlanStore
	clock OpsClock
	seq   uint64
}

func newPlanService(store *PlanStore) *PlanService {
	return &PlanService{store: store, clock: newOpsClock()}
}

func (s *PlanService) Create(ctx context.Context, p Plan) (Plan, error) {
	select {
	case <-ctx.Done():
		return Plan{}, ctx.Err()
	default:
	}
	p = normalizePlan(p)
	if strings.TrimSpace(p.WellID) == "" {
		return Plan{}, fmt.Errorf("%w: well required", ErrOpsInvalid)
	}
	if strings.TrimSpace(p.Subject) == "" {
		return Plan{}, fmt.Errorf("%w: subject required", ErrOpsInvalid)
	}
	p.ID = fmt.Sprintf("plan-%06d", atomic.AddUint64(&s.seq, 1))
	now := s.clock.Now().UTC()
	p.CreatedAt = now
	p.NextDueAt = planNextDue(now, p.IntervalDays)
	if err := s.store.Put(ctx, p); err != nil {
		return Plan{}, err
	}
	return p, nil
}

func (s *PlanService) MarkDone(ctx context.Context, id string, expected int) (Plan, error) {
	select {
	case <-ctx.Done():
		return Plan{}, ctx.Err()
	default:
	}
	p, err := s.store.Get(ctx, id)
	if err != nil {
		return Plan{}, err
	}
	if expected > 0 && expected != p.Revision {
		return Plan{}, ErrOpsConflict
	}
	if p.Status == PlanArchived {
		return Plan{}, fmt.Errorf("%w: archived plan cannot be completed", ErrOpsTransition)
	}
	now := s.clock.Now().UTC()
	p.LastDoneAt = now
	p.NextDueAt = planNextDue(now, p.IntervalDays)
	if p.Status == PlanPaused {
		p.Status = PlanActive
	}
	updated, err := s.store.Update(ctx, p, expected)
	if err != nil {
		return Plan{}, err
	}
	return updated, nil
}

func (s *PlanService) Due(ctx context.Context, now time.Time) ([]Plan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	items, err := s.store.List(ctx)
	if err != nil {
		return nil, err
	}
	out := []Plan{}
	for _, p := range items {
		if planIsActive(p) && !p.NextDueAt.After(now) {
			out = append(out, p)
		}
	}
	return out, nil
}

func (s *PlanService) List(ctx context.Context) ([]Plan, error) { return s.store.List(ctx) }

func (s *PlanService) Get(ctx context.Context, id string) (Plan, error) { return s.store.Get(ctx, id) }
