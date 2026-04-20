package watch

import (
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func aggPort(n int) scanner.Port { return scanner.Port{Number: n, Proto: "tcp"} }

func TestAggregatorAddAndFlush(t *testing.T) {
	a := NewAggregator(0)
	a.Add(aggPort(80))
	a.Add(aggPort(443))

	evts := a.Flush()
	if len(evts) != 2 {
		t.Fatalf("expected 2 events, got %d", len(evts))
	}
	if evts[0].Number != 80 || evts[1].Number != 443 {
		t.Fatalf("unexpected events: %v", evts)
	}
}

func TestAggregatorFlushClearsBuffer(t *testing.T) {
	a := NewAggregator(0)
	a.Add(aggPort(80))
	a.Flush()
	if a.Len() != 0 {
		t.Fatal("expected buffer to be empty after flush")
	}
}

func TestAggregatorFlushEmptyReturnsNil(t *testing.T) {
	a := NewAggregator(0)
	if evts := a.Flush(); len(evts) != 0 {
		t.Fatalf("expected empty slice, got %v", evts)
	}
}

func TestAggregatorRespectsMax(t *testing.T) {
	a := NewAggregator(2)
	a.Add(aggPort(80))
	a.Add(aggPort(443))
	a.Add(aggPort(8080)) // should be dropped

	if a.Len() != 2 {
		t.Fatalf("expected 2 events, got %d", a.Len())
	}
}

func TestAggregatorLen(t *testing.T) {
	a := NewAggregator(0)
	if a.Len() != 0 {
		t.Fatal("expected initial len 0")
	}
	a.Add(aggPort(80))
	if a.Len() != 1 {
		t.Fatal("expected len 1 after add")
	}
}
