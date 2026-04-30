package watch

import (
	"testing"
	"time"
)

func TestValveClosedByDefault(t *testing.T) {
	v := NewValve(time.Second)
	if v.IsOpen("k") {
		t.Fatal("expected closed by default")
	}
}

func TestValveOpenTransition(t *testing.T) {
	v := NewValve(time.Second)
	if !v.Open("k") {
		t.Fatal("expected true on first open")
	}
	if v.Open("k") {
		t.Fatal("expected false when already open")
	}
}

func TestValveIsOpenAfterOpen(t *testing.T) {
	v := NewValve(time.Second)
	v.Open("k")
	if !v.IsOpen("k") {
		t.Fatal("expected open")
	}
}

func TestValveCloseTransition(t *testing.T) {
	v := NewValve(time.Second)
	v.Open("k")
	if !v.Close("k") {
		t.Fatal("expected true on first close")
	}
	if v.Close("k") {
		t.Fatal("expected false when already closed")
	}
}

func TestValveIsClosedAfterClose(t *testing.T) {
	v := NewValve(time.Second)
	v.Open("k")
	v.Close("k")
	if v.IsOpen("k") {
		t.Fatal("expected closed after Close()")
	}
}

func TestValveWindowExpiry(t *testing.T) {
	v := NewValve(20 * time.Millisecond)
	v.Open("k")
	time.Sleep(40 * time.Millisecond)
	if v.IsOpen("k") {
		t.Fatal("expected closed after window expiry")
	}
}

func TestValveIndependentKeys(t *testing.T) {
	v := NewValve(time.Second)
	v.Open("a")
	if v.IsOpen("b") {
		t.Fatal("key b should be independent of a")
	}
}

func TestValveResetClearsState(t *testing.T) {
	v := NewValve(time.Second)
	v.Open("k")
	v.Reset()
	if v.IsOpen("k") {
		t.Fatal("expected closed after Reset")
	}
}

func TestValveDefaultsInvalidWindow(t *testing.T) {
	v := NewValve(0)
	if v.window != time.Minute {
		t.Fatalf("expected 1m default, got %v", v.window)
	}
}
