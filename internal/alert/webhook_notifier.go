package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"portwatch/internal/scanner"
)

// WebhookNotifier sends alert events to an HTTP endpoint as JSON.
type WebhookNotifier struct {
	url    string
	client *http.Client
}

type webhookPayload struct {
	Kind      string `json:"kind"`
	Proto     string `json:"proto"`
	Port      int    `json:"port"`
	Timestamp string `json:"timestamp"`
}

// NewWebhookNotifier creates a WebhookNotifier that posts to the given URL.
func NewWebhookNotifier(url string, timeout time.Duration) *WebhookNotifier {
	return &WebhookNotifier{
		url: url,
		client: &http.Client{Timeout: timeout},
	}
}

func (w *WebhookNotifier) Notify(kind string, p scanner.Port) error {
	payload := webhookPayload{
		Kind:      kind,
		Proto:     p.Proto,
		Port:      p.Number,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal: %w", err)
	}
	resp, err := w.client.Post(w.url, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d", resp.StatusCode)
	}
	return nil
}
