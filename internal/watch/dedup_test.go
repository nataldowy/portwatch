package watch

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

func makePort(proto string, number int) scanner.Port {
	return scanner.Port{Proto: proto, Number: number}
}

func TestDedupAllowsFirstAlert(t *testing.T) {
	cd := NewCooldown(5 * time.Second)
	d := NewDedup(cd)
	if !d.Allow(makePort("tcp", 8080), "new") {
		t.Fatal("expected first alert to be allowed")
	}
}

func TestDedupBlocksDuplicate(t *testing.T) {
	cd := NewCooldown(5 * time.Second)
	d := NewDedup(cd)
	d.Allow(makePort("tcp", 8080), "new")
	if d.Allow(makePort("tcp", 8080), "new") {
		t.Fatal("expected duplicate to be blocked")
	}
}

func TestDedupDistinguishesKind(t *testing.T) {
	cd := NewCooldown(5 * time.Second)
	d := NewDedup(cd)
	d.Allow(makePort("tcp", 8080), "new")
	if !d.Allow(makePort("tcp", 8080), "closed") {
		t.Fatal("expected different kind to be allowed")
	}
}

func TestDedupResetAllowsRepeat(t *testing.T) {
	cd := NewCooldown(5 * time.Second)
	d := NewDedup(cd)
	p := makePort("tcp", 443)
	d.Allow(p, "new")
	d.Reset(p, "new")
	if !d.Allow(p, "new") {
		t.Fatal("expected allow after reset")
	}
}

func TestDedupDistinguishesProto(t *testing.T) {
	cd := NewCooldown(5 * time.Second)
	d := NewDedup(cd)
	d.Allow(makePort("tcp", 53), "new")
	if !d.Allow(makePort("udp", 53), "new") {
		t.Fatal("expected different proto to be allowed")
	}
}
