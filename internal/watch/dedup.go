package watch

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Dedup suppresses repeated alerts for the same port+kind pair within a
// cooldown window, delegating to Cooldown for timing logic.
type Dedup struct {
	mu       sync.Mutex
	cooldown *Cooldown
}

// NewDedup creates a Dedup backed by the given Cooldown.
func NewDedup(cd *Cooldown) *Dedup {
	return &Dedup{cooldown: cd}
}

// Allow returns true if this port+kind combination should produce an alert.
func (d *Dedup) Allow(port scanner.Port, kind string) bool {
	key := fmt.Sprintf("%s:%d:%s", port.Proto, port.Number, kind)
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.cooldown.Allow(key)
}

// Reset clears suppression for a specific port+kind pair.
func (d *Dedup) Reset(port scanner.Port, kind string) {
	key := fmt.Sprintf("%s:%d:%s", port.Proto, port.Number, kind)
	d.mu.Lock()
	defer d.mu.Unlock()
	d.cooldown.Reset(key)
}
