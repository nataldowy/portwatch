package watch

import (
	"sync"

	"github.com/user/portwatch/internal/scanner"
)

// Aggregator batches port events and flushes them as a slice when Flush is
// called. It is safe for concurrent use.
type Aggregator struct {
	mu     sync.Mutex
	events []scanner.Port
	max    int
}

// NewAggregator returns an Aggregator that holds at most max pending events.
// If max <= 0 it is treated as unlimited.
func NewAggregator(max int) *Aggregator {
	return &Aggregator{max: max}
}

// Add appends a port event to the pending batch. If the batch is already at
// capacity the event is silently dropped.
func (a *Aggregator) Add(p scanner.Port) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.max > 0 && len(a.events) >= a.max {
		return
	}
	a.events = append(a.events, p)
}

// Flush returns all pending events and resets the internal buffer.
func (a *Aggregator) Flush() []scanner.Port {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := a.events
	a.events = nil
	return out
}

// Len returns the number of pending events.
func (a *Aggregator) Len() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return len(a.events)
}
