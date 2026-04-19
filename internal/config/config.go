package config

import (
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	Interval      time.Duration `toml:"interval"`
	PortRange     PortRange     `toml:"port_range"`
	SnapshotPath  string        `toml:"snapshot_path"`
	HistoryPath   string        `toml:"history_path"`
	Notifiers     NotifiersCfg  `toml:"notifiers"`
	AlertCooldown time.Duration `toml:"alert_cooldown"`
}

// PortRange defines the inclusive range of ports to scan.
type PortRange struct {
	From int `toml:"from"`
	To   int `toml:"to"`
}

// Default returns a Config populated with sensible defaults.
func Default() Config {
	return Config{
		Interval:      30 * time.Second,
		PortRange:     PortRange{From: 1, To: 65535},
		SnapshotPath:  "/var/lib/portwatch/snapshot.json",
		HistoryPath:   "/var/lib/portwatch/history.ndjson",
		AlertCooldown: 5 * time.Minute,
	}
}

// Load reads a TOML config file, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return cfg, err
	}
	if err := validateRange(cfg.PortRange.From, cfg.PortRange.To); err != nil {
		return cfg, err
	}
	return cfg, nil
}
