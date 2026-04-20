package watch

import (
	"testing"
)

func TestHealthCheckDefaultsHealthy(t *testing.T) {
	hc := NewHealthCheck(3)
	s := hc.Status("scanner")
	if !s.Healthy {
		t.Fatal("expected healthy for unknown key")
	}
}

func TestHealthCheckRecordSuccess(t *testing.T) {
	hc := NewHealthCheck(3)
	hc.RecordFailure("scanner", "timeout")
	hc.RecordSuccess("scanner")
	s := hc.Status("scanner")
	if !s.Healthy {
		t.Fatal("expected healthy after success")
	}
	if s.Failures != 0 {
		t.Fatalf("expected 0 failures, got %d", s.Failures)
	}
}

func TestHealthCheckBecomesUnhealthyAtThreshold(t *testing.T) {
	hc := NewHealthCheck(3)
	for i := 0; i < 3; i++ {
		hc.RecordFailure("scanner", "err")
	}
	s := hc.Status("scanner")
	if s.Healthy {
		t.Fatal("expected unhealthy after threshold failures")
	}
}

func TestHealthCheckStillHealthyBelowThreshold(t *testing.T) {
	hc := NewHealthCheck(3)
	hc.RecordFailure("scanner", "err")
	hc.RecordFailure("scanner", "err")
	s := hc.Status("scanner")
	if !s.Healthy {
		t.Fatal("expected healthy below threshold")
	}
	if s.Failures != 2 {
		t.Fatalf("expected 2 failures, got %d", s.Failures)
	}
}

func TestHealthCheckResetRestoresHealth(t *testing.T) {
	hc := NewHealthCheck(2)
	hc.RecordFailure("scanner", "err")
	hc.RecordFailure("scanner", "err")
	hc.Reset("scanner")
	s := hc.Status("scanner")
	if !s.Healthy {
		t.Fatal("expected healthy after reset")
	}
}

func TestHealthCheckDefaultThreshold(t *testing.T) {
	hc := NewHealthCheck(0) // should default to 3
	for i := 0; i < 2; i++ {
		hc.RecordFailure("k", "err")
	}
	if !hc.Status("k").Healthy {
		t.Fatal("expected healthy before default threshold")
	}
	hc.RecordFailure("k", "err")
	if hc.Status("k").Healthy {
		t.Fatal("expected unhealthy at default threshold")
	}
}
