package watch

import (
	"sync"
	"time"
)

// Budget tracks a rolling spend against a fixed allowance per time window.
// Once the allowance is exhausted the key is blocked until the window resets.
type Budget struct {
	mu        sync.Mutex
	allowance int
	window    time.Duration
	spend     map[string]int
	reset     map[string]time.Time
	now       func() time.Time
}

// NewBudget returns a Budget that permits up to allowance units per window.
func NewBudget(allowance int, window time.Duration) *Budget {
	if allowance < 1 {
		allowance = 1
	}
	if window <= 0 {
		window = time.Minute
	}
	return &Budget{
		allowance: allowance,
		window:    window,
		spend:     make(map[string]int),
		reset:     make(map[string]time.Time),
		now:       time.Now,
	}
}

// Allow attempts to spend cost units from key's budget.
// Returns true when sufficient budget remains; false otherwise.
func (b *Budget) Allow(key string, cost int) bool {
	if cost < 1 {
		cost = 1
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	if t, ok := b.reset[key]; !ok || now.After(t) {
		b.spend[key] = 0
		b.reset[key] = now.Add(b.window)
	}
	if b.spend[key]+cost > b.allowance {
		return false
	}
	b.spend[key] += cost
	return true
}

// Remaining returns the unused budget for key in the current window.
func (b *Budget) Remaining(key string) int {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := b.now()
	if t, ok := b.reset[key]; !ok || now.After(t) {
		return b.allowance
	}
	rem := b.allowance - b.spend[key]
	if rem < 0 {
		return 0
	}
	return rem
}

// Reset clears the spend record for key, restoring the full allowance.
func (b *Budget) Reset(key string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.spend, key)
	delete(b.reset, key)
}
