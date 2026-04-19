package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"portwatch/internal/history"
	"portwatch/internal/report"
)

// runReport is invoked when the user runs: portwatch report [flags]
func runReport(historyPath string, args []string) error {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	kind := fs.String("kind", "", "filter by event kind: new or closed")
	sinceStr := fs.String("since", "", "only events after this RFC3339 timestamp")
	untilStr := fs.String("until", "", "only events before this RFC3339 timestamp")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var f report.Filter
	f.Kind = *kind
	if *sinceStr != "" {
		t, err := time.Parse(time.RFC3339, *sinceStr)
		if err != nil {
			return fmt.Errorf("invalid -since: %w", err)
		}
		f.Since = t
	}
	if *untilStr != "" {
		t, err := time.Parse(time.RFC3339, *untilStr)
		if err != nil {
			return fmt.Errorf("invalid -until: %w", err)
		}
		f.Until = t
	}

	log, err := history.NewLog(historyPath)
	if err != nil {
		return fmt.Errorf("open history: %w", err)
	}
	entries, err := log.ReadAll()
	if err != nil {
		return fmt.Errorf("read history: %w", err)
	}

	filtered := report.Apply(entries, f)
	// Build a synthetic summary from filtered entries.
	var s report.Summary
	s.TotalEvents = len(filtered)
	for _, e := range filtered {
		switch e.Kind {
		case "new":
			s.NewPorts++
		case "closed":
			s.ClosedPorts++
		}
		if s.FirstSeen.IsZero() || e.Timestamp.Before(s.FirstSeen) {
			s.FirstSeen = e.Timestamp
		}
		if e.Timestamp.After(s.LastSeen) {
			s.LastSeen = e.Timestamp
		}
	}
	report.Print(os.Stdout, s)
	return nil
}
