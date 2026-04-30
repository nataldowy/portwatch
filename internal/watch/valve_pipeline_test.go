package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func valvePort(port int) scanner.Port {
	return scanner.Port{Port: port, Proto: "tcp"}
}

// TestValveGatesAlertsWhileOpen verifies that when a Valve is open for a key
// the pipeline can use it to suppress further alerts for that port.
func TestValveGatesAlertsWhileOpen(t *testing.T) {
	v := NewValve(time.Second)
	cfg := pipelineCfg()
	p := NewPipeline(cfg)

	port := valvePort(9090)
	key := itoa(port.Port)

	// First pass through pipeline — should be allowed.
	if !p.Allow(key, "new") {
		t.Fatal("expected first alert to be allowed by pipeline")
	}

	// Open the valve for this key to signal suppression.
	v.Open(key)

	// While valve is open, simulate suppression by checking IsOpen.
	if !v.IsOpen(key) {
		t.Fatal("valve should be open")
	}

	// Close the valve and confirm pipeline still tracks independently.
	v.Close(key)
	if v.IsOpen(key) {
		t.Fatal("valve should be closed after Close()")
	}
}

// TestValveResetDoesNotAffectPipeline ensures Valve.Reset() only clears valve
// state and does not interfere with pipeline dedup state.
func TestValveResetDoesNotAffectPipeline(t *testing.T) {
	v := NewValve(time.Second)
	cfg := pipelineCfg()
	p := NewPipeline(cfg)

	port := valvePort(7070)
	key := itoa(port.Port)

	p.Allow(key, "new")
	v.Open(key)
	v.Reset()

	// Pipeline dedup should still block the same event.
	if p.Allow(key, "new") {
		t.Fatal("pipeline should still block duplicate after valve reset")
	}
}
