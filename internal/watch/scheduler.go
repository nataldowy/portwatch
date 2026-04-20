package watch

import (
	"context"
	"sync"
	"time"
)

// Task is a function executed by the Scheduler on each tick.
type Task func(ctx context.Context) error

// Scheduler drives a Task on a regular interval using a Ticker.
// It records the last run time and any error for observability.
type Scheduler struct {
	ticker  *Ticker
	task    Task
	mu      sync.RWMutex
	lastRun time.Time
	lastErr error
	runs    int
}

// NewScheduler creates a Scheduler that executes task at the given interval.
func NewScheduler(interval time.Duration, jitter float64, task Task) *Scheduler {
	return &Scheduler{
		ticker: NewTicker(interval, jitter),
		task:   task,
	}
}

// Run starts the scheduler loop. It blocks until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	go s.ticker.Run(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case ts, ok := <-s.ticker.C():
			if !ok {
				return
			}
			err := s.task(ctx)
			s.mu.Lock()
			s.lastRun = ts
			s.lastErr = err
			s.runs++
			s.mu.Unlock()
		}
	}
}

// Stats returns a snapshot of scheduler execution statistics.
func (s *Scheduler) Stats() (lastRun time.Time, lastErr error, runs int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastRun, s.lastErr, s.runs
}
