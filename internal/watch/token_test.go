package watch

import (
	"testing"
	"time"
)

func TestTokenAllowsUpToCapacity(t *testing.T) {
	tk := NewToken(10, 3)
	for i := 0; i < 3; i++ {
		if !tk.Allow("k") {
			t.Fatalf("expected allow on attempt %d", i+1)
		}
	}
	if tk.Allow("k") {
		t.Fatal("expected deny after capacity exhausted")
	}
}

func TestTokenReplenishesOverTime(t *testing.T) {
	base := time.Now()
	tk := NewToken(2, 2) // 2 tokens/sec, cap 2
	tk.now = func() time.Time { return base }

	tk.Allow("k") // consume both
	tk.Allow("k")
	if tk.Allow("k") {
		t.Fatal("bucket should be empty")
	}

	// advance 1 second — should replenish 2 tokens
	tk.now = func() time.Time { return base.Add(time.Second) }
	if !tk.Allow("k") {
		t.Fatal("expected token after replenish")
	}
}

func TestTokenIndependentKeys(t *testing.T) {
	tk := NewToken(10, 1)
	if !tk.Allow("a") {
		t.Fatal("expected allow for key a")
	}
	if !tk.Allow("b") {
		t.Fatal("expected allow for key b")
	}
	if tk.Allow("a") {
		t.Fatal("expected deny for exhausted key a")
	}
}

func TestTokenRemaining(t *testing.T) {
	tk := NewToken(10, 5)
	if tk.Remaining("x") != 5 {
		t.Fatalf("expected 5 remaining for unseen key, got %d", tk.Remaining("x"))
	}
	tk.Allow("x")
	tk.Allow("x")
	if tk.Remaining("x") != 3 {
		t.Fatalf("expected 3 remaining, got %d", tk.Remaining("x"))
	}
}

func TestTokenResetRestoresBucket(t *testing.T) {
	tk := NewToken(10, 2)
	tk.Allow("k")
	tk.Allow("k")
	if tk.Allow("k") {
		t.Fatal("bucket should be empty before reset")
	}
	tk.Reset("k")
	if !tk.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestTokenDefaultsInvalidParams(t *testing.T) {
	tk := NewToken(-1, 0)
	// should not panic and should allow at least one
	if !tk.Allow("k") {
		t.Fatal("expected at least one allow with defaulted params")
	}
}
