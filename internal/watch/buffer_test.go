package watch

import (
	"testing"
)

func TestBufferAddAndFlush(t *testing.T) {
	b := NewBuffer(10)
	b.Add("80")
	b.Add("443")
	b.Add("80")
	entries := b.Flush()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Key == "80" && e.Count != 2 {
			t.Errorf("expected count 2 for 80, got %d", e.Count)
		}
		if e.Key == "443" && e.Count != 1 {
			t.Errorf("expected count 1 for 443, got %d", e.Count)
		}
	}
}

func TestBufferFlushClears(t *testing.T) {
	b := NewBuffer(10)
	b.Add("22")
	b.Flush()
	if b.Len() != 0 {
		t.Errorf("expected empty buffer after flush")
	}
}

func TestBufferDropsOldestWhenFull(t *testing.T) {
	b := NewBuffer(3)
	b.Add("1")
	b.Add("2")
	b.Add("3")
	b.Add("4") // should drop "1"
	entries := b.Flush()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e.Key == "1" {
			t.Error("oldest entry should have been dropped")
		}
	}
}

func TestBufferResetClearsEntries(t *testing.T) {
	b := NewBuffer(10)
	b.Add("80")
	b.Add("443")
	b.Reset()
	if b.Len() != 0 {
		t.Errorf("expected 0 after reset, got %d", b.Len())
	}
}

func TestBufferDefaultsInvalidMax(t *testing.T) {
	b := NewBuffer(0)
	for i := 0; i < 70; i++ {
		b.Add("port")
	}
	if b.Len() != 1 {
		t.Errorf("expected 1 distinct key, got %d", b.Len())
	}
}
