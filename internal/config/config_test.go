package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "portwatch.toml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempConfig: %v", err)
	}
	return p
}

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	if cfg.Interval != 30*time.Second {
		t.Errorf("interval: got %v, want 30s", cfg.Interval)
	}
	if cfg.PortRange.From != 1 || cfg.PortRange.To != 65535 {
		t.Errorf("port range: got %d-%d", cfg.PortRange.From, cfg.PortRange.To)
	}
	if cfg.AlertCooldown != 5*time.Minute {
		t.Errorf("alert_cooldown: got %v, want 5m", cfg.AlertCooldown)
	}
}

func TestLoadConfig(t *testing.T) {
	p := writeTempConfig(t, `
interval = "10s"
alert_cooldown = "2m"
[port_range]
  from = 1024
  to   = 9999
`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Interval != 10*time.Second {
		t.Errorf("interval: got %v", cfg.Interval)
	}
	if cfg.PortRange.From != 1024 || cfg.PortRange.To != 9999 {
		t.Errorf("port range: got %d-%d", cfg.PortRange.From, cfg.PortRange.To)
	}
	if cfg.AlertCooldown != 2*time.Minute {
		t.Errorf("alert_cooldown: got %v", cfg.AlertCooldown)
	}
}

func TestLoadConfigMissingFile(t *testing.T) {
	cfg, err := Load("/nonexistent/portwatch.toml")
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if cfg.Interval != 30*time.Second {
		t.Errorf("expected default interval, got %v", cfg.Interval)
	}
}

func TestDurationRoundTrip(t *testing.T) {
	p := writeTempConfig(t, `interval = "1m30s"`)
	cfg, err := Load(p)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.Interval != 90*time.Second {
		t.Errorf("got %v, want 90s", cfg.Interval)
	}
}

func TestLoadConfigInvalidTOML(t *testing.T) {
	p := writeTempConfig(t, `interval = [not valid toml`)
	_, err := Load(p)
	if err == nil {
		t.Fatal("expected error for invalid TOML, got nil")
	}
}
