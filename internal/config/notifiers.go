package config

// SlackConfig holds Slack webhook notification settings.
type SlackConfig struct {
	Enabled    bool   `toml:"enabled"`
	WebhookURL string `toml:"webhook_url"`
}

// NotifiersConfig aggregates all optional notifier configurations.
type NotifiersConfig struct {
	Slack   SlackConfig   `toml:"slack"`
	Webhook WebhookConfig `toml:"webhook"`
	Email   EmailConfig   `toml:"email"`
}

// WebhookConfig holds generic webhook settings.
type WebhookConfig struct {
	Enabled bool   `toml:"enabled"`
	URL     string `toml:"url"`
}

// EmailConfig holds SMTP email notification settings.
type EmailConfig struct {
	Enabled  bool   `toml:"enabled"`
	SMTPHost string `toml:"smtp_host"`
	SMTPPort int    `toml:"smtp_port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
	To       string `toml:"to"`
}

// DefaultNotifiers returns a NotifiersConfig with safe defaults.
func DefaultNotifiers() NotifiersConfig {
	return NotifiersConfig{
		Email: EmailConfig{SMTPPort: 587},
	}
}
