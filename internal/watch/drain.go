package watch

import (
	"sync"
	"time"
)

// Drain collects events over a flush interval and releases them in batches.
// It is useful for coalescing bursts of alerts before forwarding downstream.
type Drain struct {
	mu       sync.Mutex
	bucket   map[string][]time.Time
	interval time.Duration
	max      int
}

// NewDrain returns a Drain that batches events over interval, holding at most
// max entries per key. A max <= 0 defaults to 64.
func NewDrain(interval time.Duration, max int) *Drain {
	if max <= 0 {
		max = 64
	}
	if interval <= 0 {
		interval = time.Second
	}
	return &Drain{
		bucket:   make(map[string][]time.Time),
		interval: interval,
		max:      max,
	}
}

// Add records an occurrence for key at now. Returns true when the bucket has
// reached max capacity and should be flushed by the caller.
func (d *Drain) Add(key string, now time.Time) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.bucket[key] = append(d.bucket[key], now)
	return len(d.bucket[key]) >= d.max
}

// Flush returns and clears all events for key that fall within the drain
// interval ending at now. Events outside the window are discarded.
func (d *Drain) Flush(key string, now time.Time) []time.Time {
	d.mu.Lock()
	defer d.mu.Unlock()
	events, ok := d.bucket[key]
	if !ok || len(events) == 0 {
		return nil
	}
	cutoff := now.Add(-d.interval)
	var kept []time.Time
	for _, t := range events {
		if !t.Before(cutoff) {
			kept = append(kept, t)
		}
	}
	delete(d.bucket, key)
	return kept
}

// Count returns the number of buffered events for key.
func (d *Drain) Count(key string) int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.bucket[key])
}

// Reset clears all buffered events across all keys.
func (d *Drain) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.bucket = make(map[string][]time.Time)
}
