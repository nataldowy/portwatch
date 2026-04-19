package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/daemon"
	"portwatch/internal/scanner"
	"portwatch/internal/snapshot"
)

const version = "0.1.0"

func main() {
	cfgPath := flag.String("config", "", "path to config file (default: built-in defaults)")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	// Load configuration — fall back to defaults if no file is specified.
	var cfg config.Config
	var err error
	if *cfgPath != "" {
		cfg, err = config.Load(*cfgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
			os.Exit(1)
		}
	} else {
		cfg = config.Default()
	}

	// Wire up dependencies.
	sc := scanner.NewScanner()
	store := snapshot.NewStore(cfg.SnapshotPath)
	notifier := alert.NewLogNotifier(os.Stdout)
	dispatcher := alert.NewDispatcher(notifier)

	d := daemon.New(cfg, sc, store, dispatcher)

	// Handle OS signals for graceful shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Fprintf(os.Stderr, "\nreceived signal %s, shutting down...\n", sig)
		cancel()
	}()

	fmt.Printf("portwatch %s starting (interval: %s, range: %d-%d)\n",
		version, cfg.Interval, cfg.PortRange.From, cfg.PortRange.To)

	if err := d.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "daemon error: %v\n", err)
		os.Exit(1)
	}
}
