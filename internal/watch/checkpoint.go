package watch

import (
	"sync"
	"time"
)

// Checkpoint tracks the last successful scan time and sequence number,
// allowing the watcher to resume gracefully after restarts or errors.
type Checkpoint struct {
	mu       sync.RWMutex
	lastScan time.Time
	seq      uint64
	tag      string
}

// NewCheckpoint returns an initialised Checkpoint with a zero time.
func NewCheckpoint() *Checkpoint {
	return &Checkpoint{}
}

// Record marks a successful scan at the given time with an optional tag.
func (c *Checkpoint) Record(t time.Time, tag string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastScan = t
	c.seq++
	c.tag = tag
}

// LastScan returns the time of the most recent successful scan.
func (c *Checkpoint) LastScan() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.lastScan
}

// Seq returns the monotonically increasing scan sequence number.
func (c *Checkpoint) Seq() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.seq
}

// Tag returns the tag associated with the last recorded checkpoint.
func (c *Checkpoint) Tag() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tag
}

// Age returns the duration elapsed since the last successful scan.
// Returns zero if no scan has been recorded yet.
func (c *Checkpoint) Age() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.lastScan.IsZero() {
		return 0
	}
	return time.Since(c.lastScan)
}

// Reset clears all checkpoint state back to zero values.
func (c *Checkpoint) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.lastScan = time.Time{}
	c.seq = 0
	c.tag = ""
}
