package report

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"portwatch/internal/history"
)

// Summary holds aggregated stats derived from history entries.
type Summary struct {
	TotalEvents  int
	NewPorts     int
	ClosedPorts  int
	FirstSeen    time.Time
	LastSeen     time.Time
}

// Summarise reads all entries from the log and returns a Summary.
func Summarise(log *history.Log) (Summary, error) {
	entries, err := log.ReadAll()
	if err != nil {
		return Summary{}, fmt.Errorf("report: read history: %w", err)
	}
	if len(entries) == 0 {
		return Summary{}, nil
	}
	s := Summary{
		TotalEvents: len(entries),
		FirstSeen:   entries[0].Timestamp,
		LastSeen:    entries[len(entries)-1].Timestamp,
	}
	for _, e := range entries {
		switch e.Kind {
		case "new":
			s.NewPorts++
		case "closed":
			s.ClosedPorts++
		}
	}
	return s, nil
}

// Print writes a human-readable report to w (defaults to os.Stdout).
func Print(w io.Writer, s Summary) {
	if w == nil {
		w = os.Stdout
	}
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "=== portwatch report ===")
	fmt.Fprintf(tw, "Total events:\t%d\n", s.TotalEvents)
	fmt.Fprintf(tw, "New ports:\t%d\n", s.NewPorts)
	fmt.Fprintf(tw, "Closed ports:\t%d\n", s.ClosedPorts)
	if !s.FirstSeen.IsZero() {
		fmt.Fprintf(tw, "First event:\t%s\n", s.FirstSeen.Format(time.RFC3339))
		fmt.Fprintf(tw, "Last event:\t%s\n", s.LastSeen.Format(time.RFC3339))
	}
	tw.Flush()
}
