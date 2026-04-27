package watch

import (
	"testing"
	"time"
)

func TestCascadeAllowsUnregisteredKey(t *testing.T) {
	c := NewCascade(time.Second)
	if !c.Allow("port:8080") {
		t.Fatal("expected unregistered key to be allowed")
	}
}

func TestCascadeSuppressesChildAfterParentFires(t *testing.T) {
	c := NewCascade(time.Second)
	c.Register("parent", "child1", "child2")
	c.Fire("parent")

	if c.Allow("child1") {
		t.Error("expected child1 to be suppressed")
	}
	if c.Allow("child2") {
		t.Error("expected child2 to be suppressed")
	}
}

func TestCascadeAllowsChildAfterWindowExpires(t *testing.T) {
	c := NewCascade(50 * time.Millisecond)
	c.Register("parent", "child")
	c.Fire("parent")

	if c.Allow("child") {
		t.Fatal("expected child to be suppressed immediately after fire")
	}

	time.Sleep(80 * time.Millisecond)
	if !c.Allow("child") {
		t.Error("expected child to be allowed after window expires")
	}
}

func TestCascadeIndependentParents(t *testing.T) {
	c := NewCascade(time.Second)
	c.Register("parentA", "childA")
	c.Register("parentB", "childB")
	c.Fire("parentA")

	if c.Allow("childA") {
		t.Error("expected childA to be suppressed")
	}
	if !c.Allow("childB") {
		t.Error("expected childB to be allowed; parentB did not fire")
	}
}

func TestCascadeResetClearsState(t *testing.T) {
	c := NewCascade(time.Second)
	c.Register("parent", "child")
	c.Fire("parent")
	c.Reset()

	if !c.Allow("child") {
		t.Error("expected child to be allowed after reset")
	}
}

func TestCascadeDefaultsInvalidWindow(t *testing.T) {
	c := NewCascade(-1)
	if c.window != 5*time.Second {
		t.Errorf("expected default window 5s, got %v", c.window)
	}
}
