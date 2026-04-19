package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds the portwatch daemon configuration.
type Config struct {
	// PortRanges is a list of [start, end] port ranges to scan.
	PortRanges [][2]int `json:"port_ranges"`
	// Interval is the duration between scans.
	Interval Duration `json:"interval"`
	// LogFile is the path to the alert log file. Empty means stdout.
	LogFile string `json:"log_file"`
}

// Duration wraps time.Duration for JSON unmarshalling from a string.
type Duration struct {
	time.Duration
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

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

// Default returns a Config with sensible defaults.
func Default() *Config {
	return &Config{
		PortRanges: [][2]int{{1, 1024}},
		Interval:   Duration{30 * time.Second},
		LogFile:    "",
	}
}

// Load reads a JSON config file from the given path.
func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	cfg := Default()
	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
