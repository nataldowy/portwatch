package watch

import (
	"context"
	"errors"
	"log"
	"time"
)

// Task is a function that can be supervised and restarted on failure.
type Task func(ctx context.Context) error

// SupervisorConfig controls restart behaviour.
type SupervisorConfig struct {
	MaxRetries int
	Delay      time.Duration
	Name       string
}

// Supervisor runs a Task and restarts it on error up to MaxRetries times.
type Supervisor struct {
	cfg   SupervisorConfig
	retry *Retry
}

// NewSupervisor creates a Supervisor with the given config.
func NewSupervisor(cfg SupervisorConfig) *Supervisor {
	return &Supervisor{
		cfg:   cfg,
		retry: NewRetry(cfg.MaxRetries, cfg.Delay),
	}
}

// isContextError reports whether err is a context cancellation or deadline error.
func isContextError(err error) bool {
	return errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
}

// Run starts the task and supervises it, restarting on non-context errors.
func (s *Supervisor) Run(ctx context.Context, task Task) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		err := task(ctx)
		if err == nil || isContextError(err) {
			return err
		}
		if !s.retry.Allow(s.cfg.Name) {
			log.Printf("supervisor[%s]: max retries reached, giving up: %v", s.cfg.Name, err)
			return err
		}
		log.Printf("supervisor[%s]: task failed (%v), restarting in %s (%d remaining)",
			s.cfg.Name, err, s.cfg.Delay, s.retry.Remaining(s.cfg.Name))
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.cfg.Delay):
		}
	}
}
