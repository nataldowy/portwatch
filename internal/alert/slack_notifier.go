package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/scanner"
)

// SlackNotifier sends alerts to a Slack incoming webhook.
type SlackNotifier struct {
	webhookURL string
	client     *http.Client
}

type slackPayload struct {
	Text string `json:"text"`
}

// NewSlackNotifier creates a SlackNotifier for the given webhook URL.
func NewSlackNotifier(webhookURL string) *SlackNotifier {
	return &SlackNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *SlackNotifier) Notify(event string, port scanner.Port) error {
	text := fmt.Sprintf("[portwatch] %s — port %d/%s (%s)",
		event, port.Number, port.Protocol, port.Process)

	body, err := json.Marshal(slackPayload{Text: text})
	if err != nil {
		return fmt.Errorf("slack: marshal: %w", err)
	}

	resp, err := s.client.Post(s.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("slack: post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack: unexpected status %d", resp.StatusCode)
	}
	return nil
}
