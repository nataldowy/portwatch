package watch

import (
	"testing"
	"time"
)

func TestLimiterAllowsUpToRate(t *testing.T) {
	l := NewLimiter(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !l.Allow("k") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if l.Allow("k") {
		t.Fatal("expected deny after rate exceeded")
	}
}

func TestLimiterIndependentKeys(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	if !l.Allow("a") {
		t.Fatal("expected allow for key a")
	}
	if !l.Allow("b") {
		t.Fatal("expected allow for key b")
	}
	if l.Allow("a") {
		t.Fatal("expected deny for key a after limit")
	}
}

func TestLimiterResetRestoresAllowance(t *testing.T) {
	l := NewLimiter(1, time.Minute)
	l.Allow("k")
	if l.Allow("k") {
		t.Fatal("expected deny before reset")
	}
	l.Reset("k")
	if !l.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestLimiterRemaining(t *testing.T) {
	l := NewLimiter(5, time.Minute)
	if l.Remaining("k") != 5 {
		t.Fatalf("expected 5 remaining, got %d", l.Remaining("k"))
	}
	l.Allow("k")
	l.Allow("k")
	if l.Remaining("k") != 3 {
		t.Fatalf("expected 3 remaining, got %d", l.Remaining("k"))
	}
}

func TestLimiterWindowExpiry(t *testing.T) {
	base := time.Now()
	l := NewLimiter(1, time.Second)
	l.nowFn = func() time.Time { return base }
	l.Allow("k")
	if l.Allow("k") {
		t.Fatal("expected deny within window")
	}
	l.nowFn = func() time.Time { return base.Add(2 * time.Second) }
	if !l.Allow("k") {
		t.Fatal("expected allow after window expired")
	}
}
