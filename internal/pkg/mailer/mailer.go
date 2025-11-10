package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/artfoxe6/quick-gin/internal/app/core/config"
)

const defaultProvider = "sendgrid"

type templateItem struct {
	Id           string
	Subject      string
	HTMLTemplate string
	TextTemplate string
}

var Template = map[string]templateItem{
	"code": {
		Id:           "xxxx-xxxxx",
		Subject:      "Verification code",
		HTMLTemplate: `<p>Your verification code is <strong>{{.code}}</strong>.</p>`,
		TextTemplate: `Your verification code is {{.code}}.`,
	},
}

type Address struct {
	Name  string
	Email string
}

type Email struct {
	template templateItem
	from     Address
	data     map[string]any
}

type Message struct {
	From       Address
	To         Address
	Subject    string
	TemplateID string
	Data       map[string]any
	HTMLBody   string
	TextBody   string
}

type Provider interface {
	Send(ctx context.Context, msg Message) error
}

var (
	providerOnce sync.Once
	providerErr  error
	mailProvider Provider
)

func New(template templateItem, data map[string]any) Email {
	from := Address{
		Name:  strings.TrimSpace(config.Mail.FromName),
		Email: strings.TrimSpace(config.Mail.FromAddress),
	}
	if from.Name == "" {
		from.Name = "Quick Gin"
	}
	if from.Email == "" {
		from.Email = "no-reply@example.com"
	}

	cloned := make(map[string]any, len(data))
	for k, v := range data {
		cloned[k] = v
	}

	return Email{
		template: template,
		from:     from,
		data:     cloned,
	}
}

func (e Email) SendTo(ctx context.Context, name, email string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if strings.TrimSpace(email) == "" {
		return errors.New("missing recipient email")
	}

	msg := Message{
		From:       e.from,
		To:         Address{Name: name, Email: email},
		Subject:    e.template.Subject,
		TemplateID: e.template.Id,
		Data:       e.data,
	}

	var err error
	if e.template.HTMLTemplate != "" {
		msg.HTMLBody, err = renderTemplate(e.template.HTMLTemplate, e.data)
		if err != nil {
			return fmt.Errorf("render html template: %w", err)
		}
	}
	if e.template.TextTemplate != "" {
		msg.TextBody, err = renderTemplate(e.template.TextTemplate, e.data)
		if err != nil {
			return fmt.Errorf("render text template: %w", err)
		}
	}

	provider, err := resolveProvider()
	if err != nil {
		return err
	}

	return provider.Send(ctx, msg)
}

func resolveProvider() (Provider, error) {
	providerOnce.Do(func() {
		name := strings.ToLower(strings.TrimSpace(config.Mail.Provider))
		if name == "" {
			name = defaultProvider
		}

		switch name {
		case "smtp":
			mailProvider, providerErr = newSMTPProvider()
		case "sendgrid":
			mailProvider, providerErr = newSendgridProvider()
		default:
			providerErr = fmt.Errorf("unsupported mail provider %q", name)
		}
	})

	return mailProvider, providerErr
}

func renderTemplate(src string, data map[string]any) (string, error) {
	tpl, err := template.New("mail").Parse(src)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
