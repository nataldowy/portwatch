package watch

import (
	"testing"
	"time"
)

func TestStampMarkNewKeyReturnsTrue(t *testing.T) {
	s := NewStamp(100 * time.Millisecond)
	if !s.Mark("a") {
		t.Fatal("expected true for first Mark on new key")
	}
}

func TestStampMarkWithinIntervalReturnsFalse(t *testing.T) {
	s := NewStamp(500 * time.Millisecond)
	s.Mark("a")
	if s.Mark("a") {
		t.Fatal("expected false when called again within interval")
	}
}

func TestStampMarkAfterIntervalReturnsTrue(t *testing.T) {
	s := NewStamp(30 * time.Millisecond)
	s.Mark("a")
	time.Sleep(50 * time.Millisecond)
	if !s.Mark("a") {
		t.Fatal("expected true after interval has elapsed")
	}
}

func TestStampIndependentKeys(t *testing.T) {
	s := NewStamp(500 * time.Millisecond)
	s.Mark("x")
	if !s.Mark("y") {
		t.Fatal("expected true for a different key")
	}
}

func TestStampLastSeen(t *testing.T) {
	s := NewStamp(100 * time.Millisecond)
	before := time.Now()
	s.Mark("k")
	after := time.Now()
	t2, ok := s.LastSeen("k")
	if !ok {
		t.Fatal("expected LastSeen to return ok")
	}
	if t2.Before(before) || t2.After(after) {
		t.Fatalf("LastSeen %v not within [%v, %v]", t2, before, after)
	}
}

func TestStampLastSeenMissingKey(t *testing.T) {
	s := NewStamp(100 * time.Millisecond)
	_, ok := s.LastSeen("missing")
	if ok {
		t.Fatal("expected false for unseen key")
	}
}

func TestStampAgeGrowsOverTime(t *testing.T) {
	s := NewStamp(100 * time.Millisecond)
	s.Mark("a")
	time.Sleep(20 * time.Millisecond)
	age, ok := s.Age("a")
	if !ok {
		t.Fatal("expected Age to return ok")
	}
	if age < 15*time.Millisecond {
		t.Fatalf("expected age >= 15ms, got %v", age)
	}
}

func TestStampAgeUnknownKey(t *testing.T) {
	s := NewStamp(100 * time.Millisecond)
	_, ok := s.Age("nope")
	if ok {
		t.Fatal("expected false for unknown key")
	}
}

func TestStampResetClearsEntries(t *testing.T) {
	s := NewStamp(500 * time.Millisecond)
	s.Mark("a")
	s.Reset()
	if !s.Mark("a") {
		t.Fatal("expected true after reset")
	}
}

func TestStampDefaultsInvalidInterval(t *testing.T) {
	s := NewStamp(0)
	if s.interval != time.Second {
		t.Fatalf("expected default interval of 1s, got %v", s.interval)
	}
}
