package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"portwatch/internal/alert"
	"portwatch/internal/scanner"
)

func TestSlackNotifierSendsPayload(t *testing.T) {
	var received map[string]string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	n := alert.NewSlackNotifier(ts.URL)
	p := scanner.Port{Number: 8080, Protocol: "tcp", Process: "myapp"}

	if err := n.Notify("new", p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(received["text"], "8080") {
		t.Errorf("expected port 8080 in text, got: %s", received["text"])
	}
	if !strings.Contains(received["text"], "myapp") {
		t.Errorf("expected process name in text, got: %s", received["text"])
	}
}

func TestSlackNotifierNonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	n := alert.NewSlackNotifier(ts.URL)
	p := scanner.Port{Number: 443, Protocol: "tcp", Process: "nginx"}

	if err := n.Notify("new", p); err == nil {
		t.Fatal("expected error on non-OK status")
	}
}

func TestSlackNotifierBadURL(t *testing.T) {
	n := alert.NewSlackNotifier("http://127.0.0.1:0/no-server")
	p := scanner.Port{Number: 22, Protocol: "tcp", Process: "sshd"}

	if err := n.Notify("new", p); err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
