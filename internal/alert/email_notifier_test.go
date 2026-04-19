package alert

import (
	"net/smtp"
	"strings"
	"testing"

	"portwatch/internal/scanner"
)

func TestEmailNotifierSendsCorrectPayload(t *testing.T) {
	var capturedAddr string
	var capturedFrom string
	var capturedTo []string
	var capturedMsg string

	cfg := EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user",
		Password: "pass",
		From:     "alerts@example.com",
		To:       []string{"admin@example.com"},
	}

	n := &emailNotifier{
		cfg: cfg,
		send: func(addr string, _ smtp.Auth, from string, to []string, msg []byte) error {
			capturedAddr = addr
			capturedFrom = from
			capturedTo = to
			capturedMsg = string(msg)
			return nil
		},
	}

	p := scanner.Port{Number: 8080, Proto: "tcp", Addr: "0.0.0.0"}
	if err := n.Notify(EventNew, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedAddr != "smtp.example.com:587" {
		t.Errorf("expected addr smtp.example.com:587, got %s", capturedAddr)
	}
	if capturedFrom != "alerts@example.com" {
		t.Errorf("unexpected from: %s", capturedFrom)
	}
	if len(capturedTo) != 1 || capturedTo[0] != "admin@example.com" {
		t.Errorf("unexpected to: %v", capturedTo)
	}
	if !strings.Contains(capturedMsg, "8080") {
		t.Error("message should contain port number")
	}
	if !strings.Contains(capturedMsg, string(EventNew)) {
		t.Error("message should contain event type")
	}
}

func TestEmailNotifierNoAuthWhenUsernameEmpty(t *testing.T) {
	var authPassed smtp.Auth

	cfg := EmailConfig{
		Host: "localhost",
		Port: 25,
		From: "noreply@local",
		To:   []string{"ops@local"},
	}

	n := &emailNotifier{
		cfg: cfg,
		send: func(_ string, a smtp.Auth, _ string, _ []string, _ []byte) error {
			authPassed = a
			return nil
		},
	}

	p := scanner.Port{Number: 22, Proto: "tcp", Addr: "127.0.0.1"}
	if err := n.Notify(EventClosed, p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if authPassed != nil {
		t.Error("expected nil auth when username is empty")
	}
}
