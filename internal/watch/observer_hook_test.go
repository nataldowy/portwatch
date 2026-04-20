package watch

import (
	"errors"
	"testing"
)

func TestObservedWatcherScanComplete(t *testing.T) {
	obs := NewObserver(0)
	ow := NewObservedWatcher(obs)

	var got ScanCompletePayload
	obs.Subscribe(EventScanComplete, func(e ObserverEvent) {
		got = e.Payload.(ScanCompletePayload)
	})

	ow.ScanComplete(ScanCompletePayload{PortsFound: 5, NewPorts: 2, Closed: 1})

	if got.PortsFound != 5 || got.NewPorts != 2 || got.Closed != 1 {
		t.Errorf("unexpected payload: %+v", got)
	}
}

func TestObservedWatcherAlertEmitted(t *testing.T) {
	obs := NewObserver(0)
	ow := NewObservedWatcher(obs)

	var got AlertPayload
	obs.Subscribe(EventAlertEmitted, func(e ObserverEvent) {
		got = e.Payload.(AlertPayload)
	})

	ow.AlertEmitted(AlertPayload{Kind: "new", Port: 8080, Proto: "tcp"})

	if got.Kind != "new" || got.Port != 8080 || got.Proto != "tcp" {
		t.Errorf("unexpected payload: %+v", got)
	}
}

func TestObservedWatcherScanError(t *testing.T) {
	obs := NewObserver(0)
	ow := NewObservedWatcher(obs)

	sentinel := errors.New("scan failed")
	var got ErrorPayload
	obs.Subscribe(EventScanError, func(e ObserverEvent) {
		got = e.Payload.(ErrorPayload)
	})

	ow.ScanError(sentinel)

	if !errors.Is(got.Err, sentinel) {
		t.Errorf("expected sentinel error, got %v", got.Err)
	}
}

func TestObservedWatcherEventNames(t *testing.T) {
	cases := []string{
		EventScanComplete,
		EventAlertEmitted,
		EventScanError,
		EventDaemonStart,
		EventDaemonStop,
	}
	for _, name := range cases {
		if name == "" {
			t.Errorf("event name must not be empty")
		}
	}
}
