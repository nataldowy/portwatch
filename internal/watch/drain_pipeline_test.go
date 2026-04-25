package watch

import (
	"testing"
	"time"
)

// TestDrainCoalescesBurstBeforePipeline verifies that a Drain can buffer a
// burst of events and only forward them once flushed, simulating integration
// with a downstream Pipeline.
func TestDrainCoalescesBurstBeforePipeline(t *testing.T) {
	dr := NewDrain(time.Minute, 5)
	now := time.Now()

	// Simulate rapid burst for the same port.
	keys := []string{"tcp:8080", "tcp:8080", "tcp:8080"}
	for _, k := range keys {
		dr.Add(k, now)
	}

	if dr.Count("tcp:8080") != 3 {
		t.Fatalf("expected 3 buffered events before flush")
	}

	events := dr.Flush("tcp:8080", now)
	if len(events) != 3 {
		t.Fatalf("expected 3 flushed events, got %d", len(events))
	}
	if dr.Count("tcp:8080") != 0 {
		t.Fatal("expected empty buffer after flush")
	}
}

// TestDrainMaxCapacityTriggersEarlyFlush confirms that Add signals when the
// bucket is full so the caller can flush proactively.
func TestDrainMaxCapacityTriggersEarlyFlush(t *testing.T) {
	dr := NewDrain(time.Minute, 4)
	now := time.Now()

	var full bool
	for i := 0; i < 4; i++ {
		full = dr.Add("tcp:9090", now)
	}
	if !full {
		t.Fatal("expected Add to return true when max capacity reached")
	}

	events := dr.Flush("tcp:9090", now)
	if len(events) != 4 {
		t.Fatalf("expected 4 events on early flush, got %d", len(events))
	}
}
