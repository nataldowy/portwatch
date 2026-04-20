package watch

import (
	"testing"
	"time"
)

func TestDebounceAllowsFirstOccurrence(t *testing.T) {
	d := NewDebounce(100 * time.Millisecond)
	if !d.Allow("port:8080") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestDebounceBlocksWithinQuietPeriod(t *testing.T) {
	d := NewDebounce(200 * time.Millisecond)
	d.Allow("port:8080")
	if d.Allow("port:8080") {
		t.Fatal("expected second call within quiet period to be blocked")
	}
}

func TestDebounceAllowsAfterQuietPeriod(t *testing.T) {
	d := NewDebounce(50 * time.Millisecond)
	d.Allow("port:9090")
	time.Sleep(120 * time.Millisecond)
	if !d.Allow("port:9090") {
		t.Fatal("expected call after quiet period to be allowed")
	}
}

func TestDebounceIndependentKeys(t *testing.T) {
	d := NewDebounce(200 * time.Millisecond)
	if !d.Allow("port:80") {
		t.Fatal("expected port:80 to be allowed")
	}
	if !d.Allow("port:443") {
		t.Fatal("expected port:443 to be allowed independently")
	}
	if d.Allow("port:80") {
		t.Fatal("expected port:80 to be blocked on second call")
	}
}

func TestDebounceDefaultsInvalidWait(t *testing.T) {
	d := NewDebounce(0)
	if d.wait != 500*time.Millisecond {
		t.Fatalf("expected default wait 500ms, got %v", d.wait)
	}
}

func TestDebounceResetAllowsRepeat(t *testing.T) {
	d := NewDebounce(200 * time.Millisecond)
	d.Allow("port:8080")
	d.Reset()
	if !d.Allow("port:8080") {
		t.Fatal("expected Allow after Reset to return true")
	}
}

func TestDebounceResetDoesNotPanic(t *testing.T) {
	d := NewDebounce(100 * time.Millisecond)
	d.Reset() // reset on empty state should not panic
}
