package watch

import (
	"context"
	"time"
)

// Ticker wraps a time.Ticker and provides a channel-based interface
// with support for context cancellation and configurable intervals.
type Ticker struct {
	interval time.Duration
	jitter   float64
	tick     chan time.Time
	stop     chan struct{}
}

// NewTicker creates a Ticker that fires at the given interval.
// An optional jitter factor (0.0–1.0) adds randomness to each interval.
func NewTicker(interval time.Duration, jitterFactor float64) *Ticker {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	return &Ticker{
		interval: interval,
		jitter:   jitterFactor,
		tick:     make(chan time.Time, 1),
		stop:     make(chan struct{}),
	}
}

// Run starts the ticker loop. It sends the current time on C() at each
// (possibly jittered) interval until ctx is cancelled.
func (t *Ticker) Run(ctx context.Context) {
	defer close(t.tick)
	j := NewJitter(t.interval, t.jitter)
	for {
		next := j.Next("tick")
		select {
		case <-ctx.Done():
			return
		case <-time.After(next):
			select {
			case t.tick <- time.Now():
			default:
			}
		}
	}
}

// C returns the channel on which tick times are delivered.
func (t *Ticker) C() <-chan time.Time {
	return t.tick
}

// Reset clears any pending tick from the channel.
func (t *Ticker) Reset() {
	select {
	case <-t.tick:
	default:
	}
}
