package mailer

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"strings"

	"github.com/artfoxe6/quick-gin/internal/app/core/config"
	gomail "gopkg.in/gomail.v2"
)

type smtpProvider struct {
	host          string
	port          int
	username      string
	password      string
	secure        bool
	skipTLSVerify bool
}

func newSMTPProvider() (Provider, error) {
	cfg := config.MailSMTP
	if cfg.Host == "" {
		return nil, errors.New("smtp host is not configured")
	}
	if cfg.Port == 0 {
		return nil, errors.New("smtp port is not configured")
	}
	username := strings.TrimSpace(cfg.Username)
	if username == "" {
		username = config.Mail.FromAddress
	}
	return &smtpProvider{
		host:          strings.TrimSpace(cfg.Host),
		port:          cfg.Port,
		username:      username,
		password:      cfg.Password,
		secure:        cfg.Secure,
		skipTLSVerify: cfg.SkipTLSVerify,
	}, nil
}

func (s *smtpProvider) Send(_ context.Context, msg Message) error {
	m := gomail.NewMessage()
	m.SetAddressHeader("From", msg.From.Email, msg.From.Name)
	m.SetAddressHeader("To", msg.To.Email, msg.To.Name)
	m.SetHeader("Subject", msg.Subject)

	html := strings.TrimSpace(msg.HTMLBody)
	text := strings.TrimSpace(msg.TextBody)
	switch {
	case html != "" && text != "":
		m.SetBody("text/html", html)
		m.AddAlternative("text/plain", text)
	case html != "":
		m.SetBody("text/html", html)
	case text != "":
		m.SetBody("text/plain", text)
	default:
		m.SetBody("text/plain", formatKeyValueBody(msg.Data))
	}

	dialer := gomail.NewDialer(s.host, s.port, s.username, s.password)
	dialer.SSL = s.secure
	if s.skipTLSVerify {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true} // #nosec G402
	}

	return dialer.DialAndSend(m)
}

func formatKeyValueBody(data map[string]any) string {
	if len(data) == 0 {
		return "-"
	}
	var b strings.Builder
	for k, v := range data {
		b.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}
	return strings.TrimSpace(b.String())
}
