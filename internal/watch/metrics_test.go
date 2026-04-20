package watch

import (
	"testing"
)

func TestMetricsRecordScan(t *testing.T) {
	m := NewMetrics()
	m.RecordScan()
	m.RecordScan()
	snap := m.Snapshot()
	if snap.ScansTotal != 2 {
		t.Fatalf("expected 2 scans, got %d", snap.ScansTotal)
	}
	if snap.LastScanAt.IsZero() {
		t.Fatal("expected LastScanAt to be set")
	}
}

func TestMetricsRecordAlert(t *testing.T) {
	m := NewMetrics()
	m.RecordAlert()
	snap := m.Snapshot()
	if snap.AlertsTotal != 1 {
		t.Fatalf("expected 1 alert, got %d", snap.AlertsTotal)
	}
	if snap.LastAlertAt.IsZero() {
		t.Fatal("expected LastAlertAt to be set")
	}
}

func TestMetricsRecordError(t *testing.T) {
	m := NewMetrics()
	m.RecordError("connection refused")
	snap := m.Snapshot()
	if snap.ErrorsTotal != 1 {
		t.Fatalf("expected 1 error, got %d", snap.ErrorsTotal)
	}
	if snap.LastErrorMsg != "connection refused" {
		t.Fatalf("unexpected error msg: %s", snap.LastErrorMsg)
	}
	if snap.LastErrorAt.IsZero() {
		t.Fatal("expected LastErrorAt to be set")
	}
}

func TestMetricsReset(t *testing.T) {
	m := NewMetrics()
	m.RecordScan()
	m.RecordAlert()
	m.RecordError("oops")
	m.Reset()
	snap := m.Snapshot()
	if snap.ScansTotal != 0 || snap.AlertsTotal != 0 || snap.ErrorsTotal != 0 {
		t.Fatal("expected all counters to be zero after reset")
	}
	if snap.LastErrorMsg != "" {
		t.Fatal("expected LastErrorMsg to be cleared")
	}
}

func TestMetricsSnapshotIsIndependent(t *testing.T) {
	m := NewMetrics()
	m.RecordScan()
	snap := m.Snapshot()
	m.RecordScan()
	if snap.ScansTotal != 1 {
		t.Fatal("snapshot should not reflect later mutations")
	}
}
