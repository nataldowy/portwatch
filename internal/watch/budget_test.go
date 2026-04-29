package watch

import (
	"testing"
	"time"
)

func TestBudgetAllowsUpToAllowance(t *testing.T) {
	b := NewBudget(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !b.Allow("k", 1) {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
	if b.Allow("k", 1) {
		t.Fatal("expected block after allowance exhausted")
	}
}

func TestBudgetIndependentKeys(t *testing.T) {
	b := NewBudget(1, time.Minute)
	if !b.Allow("a", 1) {
		t.Fatal("expected allow for key a")
	}
	if !b.Allow("b", 1) {
		t.Fatal("expected allow for key b")
	}
	if b.Allow("a", 1) {
		t.Fatal("expected block for key a")
	}
}

func TestBudgetWindowExpiry(t *testing.T) {
	now := time.Now()
	b := NewBudget(1, 50*time.Millisecond)
	b.now = func() time.Time { return now }
	b.Allow("k", 1)
	if b.Allow("k", 1) {
		t.Fatal("expected block within window")
	}
	b.now = func() time.Time { return now.Add(60 * time.Millisecond) }
	if !b.Allow("k", 1) {
		t.Fatal("expected allow after window expiry")
	}
}

func TestBudgetRemaining(t *testing.T) {
	b := NewBudget(5, time.Minute)
	if b.Remaining("k") != 5 {
		t.Fatalf("expected 5, got %d", b.Remaining("k"))
	}
	b.Allow("k", 2)
	if b.Remaining("k") != 3 {
		t.Fatalf("expected 3, got %d", b.Remaining("k"))
	}
}

func TestBudgetResetAllowsRepeat(t *testing.T) {
	b := NewBudget(1, time.Minute)
	b.Allow("k", 1)
	if b.Allow("k", 1) {
		t.Fatal("expected block before reset")
	}
	b.Reset("k")
	if !b.Allow("k", 1) {
		t.Fatal("expected allow after reset")
	}
}

func TestBudgetCostGreaterThanOne(t *testing.T) {
	b := NewBudget(10, time.Minute)
	if !b.Allow("k", 7) {
		t.Fatal("expected allow for cost 7")
	}
	if b.Allow("k", 4) {
		t.Fatal("expected block: 7+4 > 10")
	}
	if !b.Allow("k", 3) {
		t.Fatal("expected allow for cost 3: 7+3 == 10")
	}
}
