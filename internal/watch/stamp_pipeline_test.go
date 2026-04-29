package watch

import (
	"testing"
	"time"

	"github.com/netwatch/portwatch/internal/scanner"
)

// TestStampFiltersRepeatAlertsInPipeline verifies that a Stamp used
// before the pipeline suppresses duplicate events within the interval.
func TestStampFiltersRepeatAlertsInPipeline(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	stamp := NewStamp(200 * time.Millisecond)

	port := scanner.Port{Port: 9100, Proto: "tcp"}

	filterAndSend := func() bool {
		key := itoa(port.Port)
		if !stamp.Mark(key) {
			return false
		}
		return pl.Allow("new", port)
	}

	if !filterAndSend() {
		t.Fatal("first call should be allowed")
	}
	if filterAndSend() {
		t.Fatal("second call within interval should be suppressed by stamp")
	}
}

// TestStampAllowsAfterIntervalThenPipelineDecides verifies that once
// the Stamp interval elapses, control passes back to the pipeline.
func TestStampAllowsAfterIntervalThenPipelineDecides(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	stamp := NewStamp(30 * time.Millisecond)

	port := scanner.Port{Port: 9200, Proto: "tcp"}
	key := itoa(port.Port)

	stamp.Mark(key)
	pl.Allow("new", port)

	time.Sleep(50 * time.Millisecond)

	// After the stamp interval the pipeline dedup window is also reset
	pl.Reset()

	if !stamp.Mark(key) {
		t.Fatal("stamp should allow after interval")
	}
	if !pl.Allow("new", port) {
		t.Fatal("pipeline should allow after reset")
	}
}
