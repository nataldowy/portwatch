package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

// Config holds all portwatch runtime configuration.
type Config struct {
	PortRange  PortRange     `json:"port_range"`
	Interval   Duration      `json:"interval"`
	DataDir    string        `json:"data_dir"`
	WebhookURL string        `json:"webhook_url,omitempty"`
	WebhookTimeout Duration  `json:"webhook_timeout,omitempty"`
}

type PortRange struct {
	From int `json:"from"`
	To   int `json:"to"`
}

// Duration wraps time.Duration for JSON marshalling as a string.
type Duration struct{ time.Duration }

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	parsed, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = parsed
	return nil
}

// Default returns a Config with sensible defaults.
func Default() Config {
	return Config{
		PortRange:      PortRange{From: 1, To: 65535},
		Interval:       Duration{30 * time.Second},
		DataDir:        "/var/lib/portwatch",
		WebhookTimeout: Duration{5 * time.Second},
	}
}

// Load reads a JSON config file, falling back to defaults for missing fields.
func Load(path string) (Config, error) {
	cfg := Default()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, fmt.Errorf("config file not found: %s", path)
		}
		return cfg, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parse config: %w", err)
	}
	if err := validateRange(cfg.PortRange.From, cfg.PortRange.To); err != nil {
		return cfg, err
	}
	return cfg, nil
}
