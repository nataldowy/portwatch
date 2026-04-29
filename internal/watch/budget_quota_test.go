package watch

import (
	"testing"
	"time"
)

// TestBudgetAndQuotaComposed verifies that composing Budget (coarse daily cap)
// with Quota (per-minute fine-grained limit) produces the tighter of the two
// constraints.
func TestBudgetAndQuotaComposed(t *testing.T) {
	// Budget: 5 per hour; Quota: 2 per minute.
	b := NewBudget(5, time.Hour)
	q := NewQuota(2, time.Minute)

	pass := func(key string) bool {
		if !q.Allow(key) {
			return false
		}
		if !b.Allow(key, 1) {
			return false
		}
		return true
	}

	// First two calls should pass both guards.
	if !pass("p") || !pass("p") {
		t.Fatal("expected first two calls to pass")
	}
	// Third call blocked by quota (2 per minute exhausted).
	if pass("p") {
		t.Fatal("expected third call to be blocked by quota")
	}
	// Budget should show 3 remaining (only 2 were charged).
	if rem := b.Remaining("p"); rem != 3 {
		t.Fatalf("expected budget remaining 3, got %d", rem)
	}
}

// TestBudgetExhaustedBeforeQuota verifies that when the budget runs out first
// the quota still has capacity but no events pass.
func TestBudgetExhaustedBeforeQuota(t *testing.T) {
	b := NewBudget(2, time.Hour)
	q := NewQuota(10, time.Minute)

	passed := 0
	for i := 0; i < 5; i++ {
		if q.Allow("x") && b.Allow("x", 1) {
			passed++
		}
	}
	if passed != 2 {
		t.Fatalf("expected 2 passed (budget cap), got %d", passed)
	}
	// Quota still has headroom.
	if rem := q.Remaining("x"); rem != 8 {
		t.Fatalf("expected quota remaining 8, got %d", rem)
	}
}
