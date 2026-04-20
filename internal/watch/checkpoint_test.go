package watch

import (
	"testing"
	"time"
)

func TestCheckpointDefaultsZero(t *testing.T) {
	cp := NewCheckpoint()
	if !cp.LastScan().IsZero() {
		t.Errorf("expected zero time, got %v", cp.LastScan())
	}
	if cp.Seq() != 0 {
		t.Errorf("expected seq 0, got %d", cp.Seq())
	}
	if cp.Tag() != "" {
		t.Errorf("expected empty tag, got %q", cp.Tag())
	}
	if cp.Age() != 0 {
		t.Errorf("expected zero age before first record, got %v", cp.Age())
	}
}

func TestCheckpointRecordUpdatesFields(t *testing.T) {
	cp := NewCheckpoint()
	now := time.Now()
	cp.Record(now, "scan-1")

	if !cp.LastScan().Equal(now) {
		t.Errorf("expected %v, got %v", now, cp.LastScan())
	}
	if cp.Seq() != 1 {
		t.Errorf("expected seq 1, got %d", cp.Seq())
	}
	if cp.Tag() != "scan-1" {
		t.Errorf("expected tag 'scan-1', got %q", cp.Tag())
	}
}

func TestCheckpointSeqIncrementsOnEachRecord(t *testing.T) {
	cp := NewCheckpoint()
	for i := uint64(1); i <= 5; i++ {
		cp.Record(time.Now(), "")
		if cp.Seq() != i {
			t.Errorf("expected seq %d, got %d", i, cp.Seq())
		}
	}
}

func TestCheckpointAgeGrowsOverTime(t *testing.T) {
	cp := NewCheckpoint()
	cp.Record(time.Now().Add(-100*time.Millisecond), "")
	age := cp.Age()
	if age < 100*time.Millisecond {
		t.Errorf("expected age >= 100ms, got %v", age)
	}
}

func TestCheckpointResetClearsState(t *testing.T) {
	cp := NewCheckpoint()
	cp.Record(time.Now(), "before-reset")
	cp.Reset()

	if !cp.LastScan().IsZero() {
		t.Errorf("expected zero time after reset, got %v", cp.LastScan())
	}
	if cp.Seq() != 0 {
		t.Errorf("expected seq 0 after reset, got %d", cp.Seq())
	}
	if cp.Tag() != "" {
		t.Errorf("expected empty tag after reset, got %q", cp.Tag())
	}
	if cp.Age() != 0 {
		t.Errorf("expected zero age after reset, got %v", cp.Age())
	}
}

func TestCheckpointTagOverwrittenOnRecord(t *testing.T) {
	cp := NewCheckpoint()
	cp.Record(time.Now(), "first")
	cp.Record(time.Now(), "second")
	if cp.Tag() != "second" {
		t.Errorf("expected tag 'second', got %q", cp.Tag())
	}
}
