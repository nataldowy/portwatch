package watch

import (
	"testing"
)

func TestRegistryRegisterAndLookup(t *testing.T) {
	r := NewRegistry()

	if err := r.Register("alpha", 42); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, ok := r.Lookup("alpha")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v.(int) != 42 {
		t.Fatalf("expected 42, got %v", v)
	}
}

func TestRegistryDuplicateKeyReturnsError(t *testing.T) {
	r := NewRegistry()

	_ = r.Register("dup", "first")
	err := r.Register("dup", "second")
	if err == nil {
		t.Fatal("expected error for duplicate key, got nil")
	}
}

func TestRegistryLookupMissingKey(t *testing.T) {
	r := NewRegistry()

	_, ok := r.Lookup("missing")
	if ok {
		t.Fatal("expected false for missing key")
	}
}

func TestRegistryUnregister(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("gone", true)
	r.Unregister("gone")

	_, ok := r.Lookup("gone")
	if ok {
		t.Fatal("expected key to be removed")
	}
}

func TestRegistryLen(t *testing.T) {
	r := NewRegistry()
	if r.Len() != 0 {
		t.Fatalf("expected 0, got %d", r.Len())
	}
	_ = r.Register("a", 1)
	_ = r.Register("b", 2)
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}

func TestRegistryKeys(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("x", nil)
	_ = r.Register("y", nil)

	keys := r.Keys()
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	seen := make(map[string]bool)
	for _, k := range keys {
		seen[k] = true
	}
	if !seen["x"] || !seen["y"] {
		t.Fatalf("unexpected keys: %v", keys)
	}
}

func TestRegistryReset(t *testing.T) {
	r := NewRegistry()
	_ = r.Register("keep", 1)
	r.Reset()

	if r.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", r.Len())
	}
}
