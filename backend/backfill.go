package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type BackfillStatus string

const (
	BackfillQueued  BackfillStatus = "queued"
	BackfillRunning BackfillStatus = "running"
	BackfillDone    BackfillStatus = "done"
	BackfillFailed  BackfillStatus = "failed"
)

type BackfillJob struct {
	ID         string
	Status     BackfillStatus
	Total      int
	Inserted   int
	StartedAt  time.Time
	FinishedAt time.Time
	ErrMsg     string

	done     chan struct{}
	finished chan struct{}
	cancel   context.CancelFunc
}

type BackfillRunner struct {
	store *ReadingStore
	clock OpsClock
	mu    sync.RWMutex
	jobs  map[string]*BackfillJob
	seq   uint64
}

func newBackfillRunner(store *ReadingStore) *BackfillRunner {
	return &BackfillRunner{store: store, clock: newOpsClock(), jobs: map[string]*BackfillJob{}}
}

func parseBackfillLine(line string) (Reading, error) {
	fields := strings.Split(line, ",")
	if len(fields) != 3 {
		return Reading{}, fmt.Errorf("%w: expected well_id,pressure,temperature", ErrOpsInvalid)
	}
	pressure, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
	if err != nil {
		return Reading{}, fmt.Errorf("%w: bad pressure %q", ErrOpsInvalid, strings.TrimSpace(fields[1]))
	}
	temp, err := strconv.ParseFloat(strings.TrimSpace(fields[2]), 64)
	if err != nil {
		return Reading{}, fmt.Errorf("%w: bad temperature %q", ErrOpsInvalid, strings.TrimSpace(fields[2]))
	}
	return normalizeReading(Reading{WellID: fields[0], PressureKPa: pressure, TemperatureC: temp, Source: "backfill"})
}

func parseBackfillPayload(raw string) ([]Reading, error) {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	out := make([]Reading, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		reading, err := parseBackfillLine(line)
		if err != nil {
			return nil, err
		}
		out = append(out, reading)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("%w: no valid lines", ErrOpsInvalid)
	}
	return out, nil
}

func (r *BackfillRunner) Start(ctx context.Context, raw string) (*BackfillJob, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	readings, err := parseBackfillPayload(raw)
	if err != nil {
		return nil, err
	}
	id := fmt.Sprintf("bf-%06d", atomic.AddUint64(&r.seq, 1))
	runCtx, cancel := context.WithCancel(context.Background())
	job := &BackfillJob{
		ID:        id,
		Status:    BackfillQueued,
		Total:     len(readings),
		StartedAt: r.clock.Now().UTC(),
		done:      make(chan struct{}),
		finished:  make(chan struct{}),
		cancel:    cancel,
	}
	snapshot := *job
	r.mu.Lock()
	r.jobs[id] = job
	r.mu.Unlock()
	go r.run(runCtx, id, readings)
	return &snapshot, nil
}

func (r *BackfillRunner) Cancel(ctx context.Context, id string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	r.mu.RLock()
	job, ok := r.jobs[id]
	r.mu.RUnlock()
	if !ok {
		return ErrOpsNotFound
	}
	job.cancel()
	return nil
}

func (r *BackfillRunner) Get(ctx context.Context, id string) (BackfillJob, error) {
	select {
	case <-ctx.Done():
		return BackfillJob{}, ctx.Err()
	default:
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	job, ok := r.jobs[id]
	if !ok {
		return BackfillJob{}, ErrOpsNotFound
	}
	clone := *job
	clone.done = nil
	clone.finished = nil
	clone.cancel = nil
	return clone, nil
}

func (r *BackfillRunner) WaitIdle(ctx context.Context, id string) error {
	r.mu.RLock()
	job, ok := r.jobs[id]
	r.mu.RUnlock()
	if !ok {
		return ErrOpsNotFound
	}
	select {
	case <-job.finished:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *BackfillRunner) doneChan(id string) chan struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.jobs[id].done
}

func (r *BackfillRunner) finishedChan(id string) chan struct{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.jobs[id].finished
}

func (r *BackfillRunner) setStatus(id string, status BackfillStatus) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if job := r.jobs[id]; job != nil {
		job.Status = status
	}
}

func (r *BackfillRunner) bumpInserted(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if job := r.jobs[id]; job != nil {
		job.Inserted++
	}
}

func (r *BackfillRunner) finalize(id string, ctx context.Context) {
	r.mu.Lock()
	defer r.mu.Unlock()
	job := r.jobs[id]
	if job == nil {
		return
	}
	job.FinishedAt = r.clock.Now().UTC()
	if ctx.Err() != nil {
		job.Status = BackfillFailed
		job.ErrMsg = ctx.Err().Error()
	} else {
		job.Status = BackfillDone
	}
}

func (r *BackfillRunner) run(ctx context.Context, id string, readings []Reading) {
	r.setStatus(id, BackfillRunning)
	ch := make(chan Reading, 8)
	done := r.doneChan(id)

	var producers sync.WaitGroup
	var consumer sync.WaitGroup

	consumer.Add(1)
	go func() {
		defer consumer.Done()
		for reading := range ch {
			r.store.Append(reading)
			r.bumpInserted(id)
		}
		close(done)
	}()

	workers := 4
	if workers > len(readings) {
		workers = len(readings)
	}
	for w := 0; w < workers; w++ {
		producers.Add(1)
		go func(worker int) {
			defer producers.Done()
			for i := worker; i < len(readings); i += workers {
				select {
				case ch <- readings[i]:
				case <-ctx.Done():
					return
				}
			}
		}(w)
	}
	go func() {
		producers.Wait()
		close(ch)
	}()

	<-done
	r.finalize(id, ctx)
	close(r.finishedChan(id))
}
