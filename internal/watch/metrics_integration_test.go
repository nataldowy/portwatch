package watch

import (
	"sync"
	"testing"
)

func TestMetricsConcurrentRecords(t *testing.T) {
	m := NewMetrics()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines * 3)

	for i := 0; i < goroutines; i++ {
		go func() { defer wg.Done(); m.RecordScan() }()
		go func() { defer wg.Done(); m.RecordAlert() }()
		go func() { defer wg.Done(); m.RecordError("err") }()
	}
	wg.Wait()

	snap := m.Snapshot()
	if snap.ScansTotal != goroutines {
		t.Fatalf("expected %d scans, got %d", goroutines, snap.ScansTotal)
	}
	if snap.AlertsTotal != goroutines {
		t.Fatalf("expected %d alerts, got %d", goroutines, snap.AlertsTotal)
	}
	if snap.ErrorsTotal != goroutines {
		t.Fatalf("expected %d errors, got %d", goroutines, snap.ErrorsTotal)
	}
}

func TestMetricsConcurrentSnapshotAndReset(t *testing.T) {
	m := NewMetrics()
	var wg sync.WaitGroup
	const workers = 20
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				m.RecordScan()
				_ = m.Snapshot()
			} else {
				m.Reset()
			}
		}(i)
	}
	wg.Wait() // should not race
}
