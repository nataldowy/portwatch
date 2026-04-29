package watch

import (
	"sync"
	"time"
)

// Spike detects sudden bursts of activity by comparing the current event
// rate against a rolling baseline. If the ratio exceeds the configured
// multiplier within the observation window, Allow returns true.
type Spike struct {
	mu         sync.Mutex
	window     time.Duration
	multiplier float64
	events     map[string][]time.Time
	baseline   map[string]float64
}

// NewSpike creates a Spike detector. window controls the observation
// period; multiplier is the factor above the baseline that triggers a
// spike (minimum 1.5, defaults to 2.0 if invalid).
func NewSpike(window time.Duration, multiplier float64) *Spike {
	if window <= 0 {
		window = 10 * time.Second
	}
	if multiplier < 1.5 {
		multiplier = 2.0
	}
	return &Spike{
		window:     window,
		multiplier: multiplier,
		events:     make(map[string][]time.Time),
		baseline:   make(map[string]float64),
	}
}

// Allow records a new event for key and returns true if the current
// rate exceeds multiplier × baseline rate, indicating a spike.
func (s *Spike) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-s.window)

	// Prune stale events.
	fresh := s.events[key][:0]
	for _, t := range s.events[key] {
		if t.After(cutoff) {
			fresh = append(fresh, t)
		}
	}
	fresh = append(fresh, now)
	s.events[key] = fresh

	current := float64(len(fresh))
	base, ok := s.baseline[key]
	if !ok || base == 0 {
		// Seed baseline on first call; not yet a spike.
		s.baseline[key] = current
		return false
	}

	spike := current >= s.multiplier*base
	// Update baseline with exponential moving average.
	s.baseline[key] = 0.8*base + 0.2*current
	return spike
}

// Reset clears all recorded events and baselines for key.
func (s *Spike) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, key)
	delete(s.baseline, key)
}

// ResetAll clears state for every key.
func (s *Spike) ResetAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = make(map[string][]time.Time)
	s.baseline = make(map[string]float64)
}
