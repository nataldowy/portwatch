package watch

import (
	"testing"
	"time"
)

func TestSignalUnfiredByDefault(t *testing.T) {
	s := NewSignal()
	if s.Fired() {
		t.Fatal("expected signal to be unfired initially")
	}
	if !s.FiredAt().IsZero() {
		t.Fatal("expected FiredAt to be zero before firing")
	}
}

func TestSignalFireUnblocksWait(t *testing.T) {
	s := NewSignal()
	done := make(chan struct{})
	go func() {
		<-s.Wait()
		close(done)
	}()
	s.Fire()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait did not unblock after Fire")
	}
}

func TestSignalFiredReportsTrue(t *testing.T) {
	s := NewSignal()
	s.Fire()
	if !s.Fired() {
		t.Fatal("expected Fired() to return true after Fire")
	}
}

func TestSignalFiredAtRecorded(t *testing.T) {
	s := NewSignal()
	before := time.Now()
	s.Fire()
	after := time.Now()
	at := s.FiredAt()
	if at.Before(before) || at.After(after) {
		t.Fatalf("FiredAt %v not in [%v, %v]", at, before, after)
	}
}

func TestSignalFireIdempotent(t *testing.T) {
	s := NewSignal()
	s.Fire()
	t1 := s.FiredAt()
	time.Sleep(2 * time.Millisecond)
	s.Fire() // second call must be a no-op
	if !s.FiredAt().Equal(t1) {
		t.Fatal("expected FiredAt to remain unchanged on repeated Fire")
	}
}

func TestSignalResetRearms(t *testing.T) {
	s := NewSignal()
	s.Fire()
	s.Reset()
	if s.Fired() {
		t.Fatal("expected Fired() to be false after Reset")
	}
	if !s.FiredAt().IsZero() {
		t.Fatal("expected FiredAt to be zero after Reset")
	}
}

func TestSignalWaitAfterFireReturnsImmediately(t *testing.T) {
	s := NewSignal()
	s.Fire()
	select {
	case <-s.Wait():
	default:
		t.Fatal("Wait on already-fired signal should not block")
	}
}

func TestSignalResetThenFireAgain(t *testing.T) {
	s := NewSignal()
	s.Fire()
	s.Reset()
	done := make(chan struct{})
	go func() {
		<-s.Wait()
		close(done)
	}()
	s.Fire()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Wait did not unblock after second Fire")
	}
}
