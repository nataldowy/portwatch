package alert

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"portwatch/internal/scanner"
)

func TestWebhookNotifierSendsPayload(t *testing.T) {
	var received webhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, 5*time.Second)
	p := scanner.Port{Number: 8080, Proto: "tcp"}
	if err := n.Notify("new", p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if received.Port != 8080 {
		t.Errorf("expected port 8080, got %d", received.Port)
	}
	if received.Proto != "tcp" {
		t.Errorf("expected proto tcp, got %s", received.Proto)
	}
	if received.Kind != "new" {
		t.Errorf("expected kind new, got %s", received.Kind)
	}
	if received.Timestamp == "" {
		t.Error("expected non-empty timestamp")
	}
}

func TestWebhookNotifierNonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	n := NewWebhookNotifier(ts.URL, 5*time.Second)
	p := scanner.Port{Number: 443, Proto: "tcp"}
	if err := n.Notify("closed", p); err == nil {
		t.Error("expected error for non-2xx status")
	}
}

func TestWebhookNotifierBadURL(t *testing.T) {
	n := NewWebhookNotifier("http://127.0.0.1:1", 1*time.Second)
	p := scanner.Port{Number: 22, Proto: "tcp"}
	if err := n.Notify("new", p); err == nil {
		t.Error("expected connection error")
	}
}
