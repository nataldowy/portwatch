package watch

import (
	"sync"
	"time"
)

// Pressure tracks the rate of events over a sliding window and reports
// whether the system is under high load (above a configurable threshold).
type Pressure struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	events    []time.Time
}

// NewPressure creates a Pressure monitor. threshold is the maximum number
// of events allowed within window before the system is considered under
// pressure. A non-positive threshold defaults to 10; an invalid window
// defaults to one minute.
func NewPressure(threshold int, window time.Duration) *Pressure {
	if threshold <= 0 {
		threshold = 10
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Pressure{threshold: threshold, window: window}
}

// Record registers a new event at the current time.
func (p *Pressure) Record() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prune(time.Now())
	p.events = append(p.events, time.Now())
}

// High returns true when the number of events recorded within the sliding
// window meets or exceeds the threshold.
func (p *Pressure) High() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prune(time.Now())
	return len(p.events) >= p.threshold
}

// Count returns the number of events currently within the sliding window.
func (p *Pressure) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.prune(time.Now())
	return len(p.events)
}

// Reset clears all recorded events.
func (p *Pressure) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.events = p.events[:0]
}

// prune removes events older than the window. Must be called with mu held.
func (p *Pressure) prune(now time.Time) {
	cutoff := now.Add(-p.window)
	i := 0
	for i < len(p.events) && p.events[i].Before(cutoff) {
		i++
	}
	p.events = p.events[i:]
}
