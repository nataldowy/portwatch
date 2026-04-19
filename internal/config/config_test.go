package config_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"portwatch/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "portwatch-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.Default()
	if len(cfg.PortRanges) == 0 {
		t.Fatal("expected default port ranges")
	}
	if cfg.Interval.Duration <= 0 {
		t.Fatal("expected positive default interval")
	}
}

func TestLoadConfig(t *testing.T) {
	path := writeTempConfig(t, `{"port_ranges":[[80,90],[443,443]],"interval":"1m","log_file":"/tmp/pw.log"}`)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if len(cfg.PortRanges) != 2 {
		t.Errorf("expected 2 ranges, got %d", len(cfg.PortRanges))
	}
	if cfg.Interval.Duration != time.Minute {
		t.Errorf("expected 1m, got %v", cfg.Interval.Duration)
	}
	if cfg.LogFile != "/tmp/pw.log" {
		t.Errorf("unexpected log_file: %s", cfg.LogFile)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/portwatch.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestDurationRoundTrip(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = config.Duration{Duration: 5 * time.Minute}
	b, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out config.Config
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if out.Interval.Duration != 5*time.Minute {
		t.Errorf("round-trip mismatch: %v", out.Interval.Duration)
	}
}
