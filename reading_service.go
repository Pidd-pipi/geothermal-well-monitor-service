package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type ReadingStats struct {
	WellID      string
	Count       int
	AvgPressure float64
	MinPressure float64
	MaxPressure float64
	AvgTemp     float64
	MinTemp     float64
	MaxTemp     float64
	FirstAt     time.Time
	LastAt      time.Time
}

type ReadingService struct {
	store *ReadingStore
	clock OpsClock

	mu          sync.RWMutex
	recentCache map[string][]Reading
	statsCache  map[string]ReadingStats
	seq         uint64
}

func newReadingService(store *ReadingStore) *ReadingService {
	return &ReadingService{
		store:       store,
		clock:       newOpsClock(),
		recentCache: map[string][]Reading{},
		statsCache:  map[string]ReadingStats{},
	}
}

func (s *ReadingService) Ingest(ctx context.Context, r Reading) (Reading, error) {
	select {
	case <-ctx.Done():
		return Reading{}, ctx.Err()
	default:
	}
	normalized, err := normalizeReading(r)
	if err != nil {
		return Reading{}, err
	}
	normalized.ID = fmt.Sprintf("rdg-%06d", atomic.AddUint64(&s.seq, 1))
	s.store.Append(normalized)

	s.mu.Lock()
	s.recentCache[normalized.WellID] = s.store.Recent(normalized.WellID, 10)
	delete(s.statsCache, normalized.WellID)
	s.mu.Unlock()
	return normalized, nil
}

func (s *ReadingService) Recent(ctx context.Context, wellID string, n int) ([]Reading, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if strings.TrimSpace(wellID) == "" {
		return nil, fmt.Errorf("%w: well_id required", ErrOpsInvalid)
	}
	if n < 1 {
		n = 10
	}
	s.mu.RLock()
	cached, ok := s.recentCache[wellID]
	s.mu.RUnlock()
	var items []Reading
	if ok {
		items = cached
	} else {
		items = s.store.Recent(wellID, n)
		s.mu.Lock()
		s.recentCache[wellID] = items
		s.mu.Unlock()
	}
	if n > len(items) {
		n = len(items)
	}
	out := make([]Reading, 0, n)
	for i := len(items) - n; i < len(items); i++ {
		out = append(out, items[i])
	}
	return out, nil
}

func (s *ReadingService) Stats(ctx context.Context, wellID string, since time.Time) (ReadingStats, error) {
	select {
	case <-ctx.Done():
		return ReadingStats{}, ctx.Err()
	default:
	}
	if strings.TrimSpace(wellID) == "" {
		return ReadingStats{}, fmt.Errorf("%w: well_id required", ErrOpsInvalid)
	}
	s.mu.RLock()
	cached, ok := s.statsCache[wellID]
	s.mu.RUnlock()
	if ok {
		return cached, nil
	}
	if since.IsZero() {
		since = time.Now().UTC().Add(-24 * time.Hour)
	}
	items := s.store.Range(wellID, since, time.Now().UTC())
	stats := ReadingStats{WellID: wellID, MinPressure: math.MaxFloat64, MinTemp: math.MaxFloat64}
	for _, item := range items {
		stats.Count++
		stats.AvgPressure += item.PressureKPa
		stats.AvgTemp += item.TemperatureC
		if item.PressureKPa < stats.MinPressure {
			stats.MinPressure = item.PressureKPa
		}
		if item.PressureKPa > stats.MaxPressure {
			stats.MaxPressure = item.PressureKPa
		}
		if item.TemperatureC < stats.MinTemp {
			stats.MinTemp = item.TemperatureC
		}
		if item.TemperatureC > stats.MaxTemp {
			stats.MaxTemp = item.TemperatureC
		}
		if stats.FirstAt.IsZero() || item.At.Before(stats.FirstAt) {
			stats.FirstAt = item.At
		}
		if item.At.After(stats.LastAt) {
			stats.LastAt = item.At
		}
	}
	if stats.Count > 0 {
		stats.AvgPressure /= float64(stats.Count)
		stats.AvgTemp /= float64(stats.Count)
	} else {
		stats.MinPressure = 0
		stats.MinTemp = 0
	}
	s.mu.Lock()
	s.statsCache[wellID] = stats
	s.mu.Unlock()
	return stats, nil
}

func (s *ReadingService) Count(wellID string) int { return s.store.Count(wellID) }
