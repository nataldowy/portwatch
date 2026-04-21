package watch

import (
	"testing"
	"time"
)

func TestLeaseGrantsFirstAcquire(t *testing.T) {
	l := NewLease(time.Second)
	ok, exp := l.Acquire("port:8080")
	if !ok {
		t.Fatal("expected first acquire to succeed")
	}
	if exp.IsZero() {
		t.Fatal("expected non-zero expiry")
	}
}

func TestLeaseBlocksWithinTTL(t *testing.T) {
	l := NewLease(time.Second)
	l.Acquire("port:8080")
	ok, _ := l.Acquire("port:8080")
	if ok {
		t.Fatal("expected second acquire within TTL to fail")
	}
}

func TestLeaseAllowsAfterExpiry(t *testing.T) {
	now := time.Now()
	l := NewLease(time.Second)
	l.now = func() time.Time { return now }

	l.Acquire("port:8080")

	// advance past TTL
	l.now = func() time.Time { return now.Add(2 * time.Second) }

	ok, _ := l.Acquire("port:8080")
	if !ok {
		t.Fatal("expected acquire after expiry to succeed")
	}
}

func TestLeaseReleaseUnblocksKey(t *testing.T) {
	l := NewLease(time.Second)
	l.Acquire("port:9090")
	l.Release("port:9090")
	ok, _ := l.Acquire("port:9090")
	if !ok {
		t.Fatal("expected acquire after release to succeed")
	}
}

func TestLeaseActiveReflectsState(t *testing.T) {
	l := NewLease(time.Second)
	if l.Active("port:443") {
		t.Fatal("expected no active lease before acquire")
	}
	l.Acquire("port:443")
	if !l.Active("port:443") {
		t.Fatal("expected active lease after acquire")
	}
	l.Release("port:443")
	if l.Active("port:443") {
		t.Fatal("expected no active lease after release")
	}
}

func TestLeaseIndependentKeys(t *testing.T) {
	l := NewLease(time.Second)
	l.Acquire("port:80")
	ok, _ := l.Acquire("port:443")
	if !ok {
		t.Fatal("expected independent keys to not interfere")
	}
}

func TestLeaseResetClearsAll(t *testing.T) {
	l := NewLease(time.Second)
	l.Acquire("port:80")
	l.Acquire("port:443")
	l.Reset()
	if l.Active("port:80") || l.Active("port:443") {
		t.Fatal("expected all leases cleared after reset")
	}
}

func TestLeaseDefaultsInvalidTTL(t *testing.T) {
	l := NewLease(-1)
	if l.ttl != 30*time.Second {
		t.Fatalf("expected default TTL 30s, got %v", l.ttl)
	}
}
