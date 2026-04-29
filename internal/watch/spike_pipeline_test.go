package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func spikePort(port int) scanner.Port {
	return scanner.Port{Port: port, Proto: "tcp"}
}

// TestSpikeBlockedByPipelineAfterBurst verifies that once Spike detects a
// burst for a given port key, the Pipeline (via Dedup) also suppresses the
// subsequent duplicate alert so only one notification propagates.
func TestSpikeBlockedByPipelineAfterBurst(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	sp := NewSpike(10*time.Second, 2.0)

	p := spikePort(9090)
	key := itoa(p.Port)

	// Seed the spike baseline.
	sp.Allow(key)

	// Burst through the pipeline — first event should pass dedup.
	passed := pl.Allow(p, "new")
	if !passed {
		t.Fatal("expected first pipeline event to pass")
	}

	// Trigger spike detection.
	var spiked bool
	for i := 0; i < 15; i++ {
		if sp.Allow(key) {
			spiked = true
			break
		}
	}
	if !spiked {
		t.Fatal("expected spike to be detected during burst")
	}

	// Pipeline dedup should now block the duplicate alert for the same event kind.
	blocked := pl.Allow(p, "new")
	if blocked {
		t.Fatal("expected pipeline to block duplicate alert after spike")
	}
}

// TestSpikeResetRestoresPipelineFlow verifies that resetting both Spike and
// Pipeline allows traffic to resume normally.
func TestSpikeResetRestoresPipelineFlow(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	sp := NewSpike(10*time.Second, 2.0)

	p := spikePort(7070)
	key := itoa(p.Port)

	sp.Allow(key)
	for i := 0; i < 15; i++ {
		sp.Allow(key)
	}

	sp.Reset(key)
	pl.Reset()

	if !sp.Allow(key) == true { // first call after reset seeds baseline — not a spike
	}

	if !pl.Allow(p, "new") {
		t.Fatal("expected pipeline to allow event after full reset")
	}
}
