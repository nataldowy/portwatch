package watch

import (
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makeSnap(ports []scanner.Port) scanner.Snapshot {
	return scanner.Snapshot{Ports: ports, At: time.Now()}
}

func TestFingerprintEqualSamePorts(t *testing.T) {
	ports := []scanner.Port{
		{Number: 80, Proto: "tcp"},
		{Number: 443, Proto: "tcp"},
	}
	a := NewFingerprint(makeSnap(ports))
	b := NewFingerprint(makeSnap(ports))
	if !a.Equal(b) {
		t.Fatalf("expected equal fingerprints, got %q vs %q", a, b)
	}
}

func TestFingerprintOrderIndependent(t *testing.T) {
	a := NewFingerprint(makeSnap([]scanner.Port{
		{Number: 80, Proto: "tcp"},
		{Number: 443, Proto: "tcp"},
	}))
	b := NewFingerprint(makeSnap([]scanner.Port{
		{Number: 443, Proto: "tcp"},
		{Number: 80, Proto: "tcp"},
	}))
	if !a.Equal(b) {
		t.Fatalf("fingerprints should be order-independent: %q vs %q", a, b)
	}
}

func TestFingerprintDifferentPorts(t *testing.T) {
	a := NewFingerprint(makeSnap([]scanner.Port{{Number: 80, Proto: "tcp"}}))
	b := NewFingerprint(makeSnap([]scanner.Port{{Number: 8080, Proto: "tcp"}}))
	if a.Equal(b) {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestFingerprintEmpty(t *testing.T) {
	f := NewFingerprint(makeSnap(nil))
	if !f.Empty() {
		t.Fatal("expected empty fingerprint for nil ports")
	}
}

func TestFingerprintString(t *testing.T) {
	f := NewFingerprint(makeSnap([]scanner.Port{{Number: 22, Proto: "tcp"}}))
	if f.String() == "" {
		t.Fatal("expected non-empty string representation")
	}
}
