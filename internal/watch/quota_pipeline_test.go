package watch

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

// TestQuotaIntegratesWithPipeline verifies that a Quota can be used alongside
// the existing Pipeline to cap total alerts per port across scan cycles.
func TestQuotaIntegratesWithPipeline(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	q := NewQuota(2, time.Minute)

	port := scanner.Port{Port: 9090, Proto: "tcp"}
	kind := "new"

	passed := 0
	for i := 0; i < 5; i++ {
		if pl.Allow(port, kind) && q.Allow(itoa(port.Port)) {
			passed++
		}
		// reset pipeline dedup so only quota gates repeat attempts
		pl.Reset()
	}

	if passed != 2 {
		t.Fatalf("expected quota to cap at 2 passes, got %d", passed)
	}
}

// TestQuotaRemainingAfterPipelineDrop confirms Remaining is unaffected when
// the pipeline blocks the event before the quota is consulted.
func TestQuotaRemainingAfterPipelineDrop(t *testing.T) {
	cfg := pipelineCfg()
	pl := NewPipeline(cfg)
	q := NewQuota(3, time.Minute)

	port := scanner.Port{Port: 7070, Proto: "tcp"}
	kind := "new"

	// First call allowed by both pipeline and quota.
	if pl.Allow(port, kind) {
		q.Allow(itoa(port.Port))
	}
	// Second call blocked by pipeline dedup — quota should not be touched.
	if pl.Allow(port, kind) {
		q.Allow(itoa(port.Port))
	}

	if r := q.Remaining(itoa(port.Port)); r != 2 {
		t.Fatalf("expected 2 remaining after one quota use, got %d", r)
	}
}
