package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func TestThrottleGroupSuppressesDuplicateAlertsInPipeline(t *testing.T) {
	cfg := pipelineCfg()
	p := NewPipeline(cfg)
	g := NewThrottleGroup(time.Second)

	port := testPort(9090, "new")

	// First event: pipeline allows, throttle group allows.
	if !p.Allow(port) {
		t.Fatal("pipeline should allow first event")
	}
	if !g.Allow(portKey(port)) {
		t.Fatal("throttle group should allow first event")
	}

	// Second event within window: pipeline dedup blocks it first.
	if p.Allow(port) {
		t.Fatal("pipeline dedup should block duplicate")
	}
}

func TestThrottleGroupAllowsAfterPipelineReset(t *testing.T) {
	cfg := pipelineCfg()
	p := NewPipeline(cfg)
	g := NewThrottleGroup(10 * time.Millisecond)

	port := testPort(9091, "new")

	p.Allow(port)
	g.Allow(portKey(port))

	p.Reset()
	g.Reset()

	if !p.Allow(port) {
		t.Fatal("pipeline should allow after reset")
	}
	if !g.Allow(portKey(port)) {
		t.Fatal("throttle group should allow after reset")
	}
}

func portKey(p scanner.Port) string {
	return itoa(p.Number) + ":" + p.Kind
}
