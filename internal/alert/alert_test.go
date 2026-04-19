package alert

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func makePort(port int, proto string) scanner.PortInfo {
	return scanner.PortInfo{Port: port, Proto: proto, PID: 42}
}

func TestLogNotifierWritesEvent(t *testing.T) {
	var buf bytes.Buffer
	n := NewLogNotifier(&buf)

	evt := Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Level:     LevelAlert,
		Message:   "new port detected",
		Port:      makePort(8080, "tcp"),
	}

	if err := n.Notify(evt); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "ALERT") {
		t.Errorf("expected ALERT in output, got: %s", line)
	}
	if !strings.Contains(line, "8080") {
		t.Errorf("expected port 8080 in output, got: %s", line)
	}
}

func TestLogNotifierDefaultsToStdout(t *testing.T) {
	n := NewLogNotifier(nil)
	if n.Out == nil {
		t.Error("expected non-nil writer")
	}
}
