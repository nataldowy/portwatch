package main

import (
	"log"
	"os"

	"portwatch/internal/alert"
	"portwatch/internal/config"
	"portwatch/internal/history"
)

// buildNotifier constructs a MultiNotifier from the provided config,
// enabling only the notifiers that are configured and active.
func buildNotifier(cfg *config.Config, log_ *history.Log) alert.Notifier {
	multi := alert.NewMultiNotifier()

	// Always include log + history notifiers.
	multi.Add(alert.NewLogNotifier(os.Stdout))
	if log_ != nil {
		multi.Add(alert.NewHistoryNotifier(log_))
	}

	n := cfg.Notifiers

	if n.Slack.Enabled && n.Slack.WebhookURL != "" {
		log.Println("[portwatch] slack notifier enabled")
		multi.Add(alert.NewSlackNotifier(n.Slack.WebhookURL))
	}

	if n.Webhook.Enabled && n.Webhook.URL != "" {
		log.Println("[portwatch] webhook notifier enabled")
		multi.Add(alert.NewWebhookNotifier(n.Webhook.URL))
	}

	if n.Email.Enabled && n.Email.SMTPHost != "" {
		log.Println("[portwatch] email notifier enabled")
		multi.Add(alert.NewEmailNotifier(
			n.Email.SMTPHost,
			n.Email.SMTPPort,
			n.Email.Username,
			n.Email.Password,
			n.Email.From,
			n.Email.To,
		))
	}

	return multi
}
