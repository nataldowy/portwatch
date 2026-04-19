package config_test

import (
	"testing"

	"portwatch/internal/config"
)

func TestDefaultNotifiersEmailPort(t *testing.T) {
	n := config.DefaultNotifiers()
	if n.Email.SMTPPort != 587 {
		t.Errorf("expected default SMTP port 587, got %d", n.Email.SMTPPort)
	}
}

func TestDefaultNotifiersSlackDisabled(t *testing.T) {
	n := config.DefaultNotifiers()
	if n.Slack.Enabled {
		t.Error("expected Slack notifier disabled by default")
	}
}

func TestDefaultNotifiersWebhookDisabled(t *testing.T) {
	n := config.DefaultNotifiers()
	if n.Webhook.Enabled {
		t.Error("expected Webhook notifier disabled by default")
	}
}

func TestLoadConfigWithSlack(t *testing.T) {
	content := `
[notifiers]
  [notifiers.slack]
    enabled = true
    webhook_url = "https://hooks.slack.com/test"
`
	path := writeTempConfig(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if !cfg.Notifiers.Slack.Enabled {
		t.Error("expected Slack enabled")
	}
	if cfg.Notifiers.Slack.WebhookURL != "https://hooks.slack.com/test" {
		t.Errorf("unexpected webhook URL: %s", cfg.Notifiers.Slack.WebhookURL)
	}
}
