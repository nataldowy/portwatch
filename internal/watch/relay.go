package watch

import (
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Relay forwards scanner.DiffEvent values to registered handlers,
// optionally tagging each event with the wall-clock time it was relayed.
type Relay struct {
	mu       sync.RWMutex
	handlers []func(scanner.DiffEvent)
	now      func() time.Time
}

// NewRelay constructs a Relay. Pass a custom clock for testing.
func NewRelay(now func() time.Time) *Relay {
	if now == nil {
		now = time.Now
	}
	return &Relay{now: now}
}

// Subscribe registers a handler that will be called for every forwarded event.
func (r *Relay) Subscribe(fn func(scanner.DiffEvent)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = append(r.handlers, fn)
}

// Forward delivers ev to all registered handlers in registration order.
func (r *Relay) Forward(ev scanner.DiffEvent) {
	r.mu.RLock()
	handlers := make([]func(scanner.DiffEvent), len(r.handlers))
	copy(handlers, r.handlers)
	r.mu.RUnlock()

	for _, h := range handlers {
		h(ev)
	}
}

// Len returns the number of registered handlers.
func (r *Relay) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.handlers)
}

// Reset removes all registered handlers.
func (r *Relay) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers = nil
}
