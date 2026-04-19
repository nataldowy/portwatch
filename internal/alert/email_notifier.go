package alert

import (
	"fmt"
	"net/smtp"
	"strings"

	"portwatch/internal/scanner"
)

// EmailConfig holds SMTP configuration for the email notifier.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// emailNotifier sends alert notifications via email.
type emailNotifier struct {
	cfg  EmailConfig
	send func(addr string, a smtp.Auth, from string, to []string, msg []byte) error
}

// NewEmailNotifier creates a Notifier that sends emails over SMTP.
func NewEmailNotifier(cfg EmailConfig) Notifier {
	return &emailNotifier{
		cfg:  cfg,
		send: smtp.SendMail,
	}
}

func (e *emailNotifier) Notify(event Event, p scanner.Port) error {
	subject := fmt.Sprintf("[portwatch] %s port %d/%s", event, p.Number, p.Proto)
	body := fmt.Sprintf(
		"Port change detected:\r\nEvent : %s\r\nPort  : %d\r\nProto : %s\r\nAddr  : %s\r\n",
		event, p.Number, p.Proto, p.Addr,
	)
	msg := []byte(strings.Join([]string{
		"From: " + e.cfg.From,
		"To: " + strings.Join(e.cfg.To, ", "),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=utf-8",
		"",
		body,
	}, "\r\n"))

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	var auth smtp.Auth
	if e.cfg.Username != "" {
		auth = smtp.PlainAuth("", e.cfg.Username, e.cfg.Password, e.cfg.Host)
	}
	return e.send(addr, auth, e.cfg.From, e.cfg.To, msg)
}
