package alert

import (
	"fmt"
	"io"
	"os"
	"time"

	"portwatch/internal/scanner"
)

// Level represents the severity of an alert.
type Level string

const (
	LevelInfo  Level = "INFO"
	LevelWarn  Level = "WARN"
	LevelAlert Level = "ALERT"
)

// Event holds information about a single port change alert.
type Event struct {
	Timestamp time.Time
	Level     Level
	Message   string
	Port      scanner.PortInfo
}

// Notifier sends alert events to a destination.
type Notifier interface {
	Notify(event Event) error
}

// LogNotifier writes alerts as formatted lines to a writer.
type LogNotifier struct {
	Out io.Writer
}

// NewLogNotifier returns a LogNotifier writing to stdout by default.
func NewLogNotifier(out io.Writer) *LogNotifier {
	if out == nil {
		out = os.Stdout
	}
	return &LogNotifier{Out: out}
}

// Notify formats and writes the event to the configured writer.
func (l *LogNotifier) Notify(event Event) error {
	_, err := fmt.Fprintf(
		l.Out,
		"[%s] %s port=%d proto=%s pid=%d\n",
		event.Timestamp.Format(time.RFC3339),
		event.Level,
		event.Port.Port,
		event.Port.Proto,
		event.Port.PID,
	)
	return err
}
