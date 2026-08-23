package main

import (
	"context"
	"time"
)

type OpsClock struct{ NowFunc func() time.Time }

func newOpsClock() OpsClock { return OpsClock{NowFunc: time.Now} }
func (c OpsClock) Now() time.Time {
	if c.NowFunc == nil {
		return time.Now().UTC()
	}
	return c.NowFunc().UTC()
}
func (c OpsClock) Stamp() string { return c.Now().Format(time.RFC3339Nano) }
func opsContext(parent context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if parent == nil {
		parent = context.Background()
	}
	if timeout <= 0 {
		timeout = 5 * time.Second
	}
	return context.WithTimeout(parent, timeout)
}
func opsDeadline(ctx context.Context) bool {
	if ctx == nil {
		return false
	}
	_, ok := ctx.Deadline()
	return ok
}
func opsParseStamp(value string) (time.Time, error) { return time.Parse(time.RFC3339Nano, value) }
func opsBackoff(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 6 {
		attempt = 6
	}
	return time.Duration(1<<uint(attempt-1)) * 20 * time.Millisecond
}
func opsDelay(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func opsAge(now time.Time, stamp string) time.Duration {
	parsed, err := opsParseStamp(stamp)
	if err != nil || now.Before(parsed) {
		return 0
	}
	return now.Sub(parsed)
}
