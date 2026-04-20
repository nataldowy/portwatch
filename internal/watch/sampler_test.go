package watch

import (
	"testing"
	"time"
)

func TestSamplerAllowsFirstOccurrence(t *testing.T) {
	s := NewSampler(time.Second)
	if !s.Allow("port:80") {
		t.Fatal("expected first occurrence to be allowed")
	}
}

func TestSamplerBlocksWithinInterval(t *testing.T) {
	now := time.Now()
	s := NewSampler(time.Second)
	s.now = func() time.Time { return now }

	s.Allow("port:80")
	if s.Allow("port:80") {
		t.Fatal("expected second occurrence within interval to be blocked")
	}
}

func TestSamplerAllowsAfterInterval(t *testing.T) {
	now := time.Now()
	s := NewSampler(time.Second)
	s.now = func() time.Time { return now }
	s.Allow("port:80")

	s.now = func() time.Time { return now.Add(2 * time.Second) }
	if !s.Allow("port:80") {
		t.Fatal("expected occurrence after interval to be allowed")
	}
}

func TestSamplerIndependentKeys(t *testing.T) {
	s := NewSampler(time.Second)
	s.Allow("port:80")
	if !s.Allow("port:443") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestSamplerResetAllowsRepeat(t *testing.T) {
	now := time.Now()
	s := NewSampler(time.Second)
	s.now = func() time.Time { return now }
	s.Allow("port:80")
	s.Reset("port:80")
	if !s.Allow("port:80") {
		t.Fatal("expected reset key to be allowed")
	}
}

func TestSamplerLastSeen(t *testing.T) {
	now := time.Now()
	s := NewSampler(time.Second)
	s.now = func() time.Time { return now }

	_, ok := s.LastSeen("port:80")
	if ok {
		t.Fatal("expected key to be unseen before Allow")
	}

	s.Allow("port:80")
	t2, ok := s.LastSeen("port:80")
	if !ok || !t2.Equal(now) {
		t.Fatalf("expected LastSeen=%v got %v ok=%v", now, t2, ok)
	}
}

func TestSamplerResetAll(t *testing.T) {
	s := NewSampler(time.Second)
	s.Allow("port:80")
	s.Allow("port:443")
	s.ResetAll()
	if !s.Allow("port:80") || !s.Allow("port:443") {
		t.Fatal("expected all keys to be reset")
	}
}
