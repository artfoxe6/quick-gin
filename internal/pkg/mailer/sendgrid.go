package mailer

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/artfoxe6/quick-gin/internal/app/core/config"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendgridProvider struct {
	apiKey string
}

func newSendgridProvider() (Provider, error) {
	key := strings.TrimSpace(config.MailSendgrid.Key)
	if key == "" {
		return nil, errors.New("sendgrid api key is not configured")
	}
	return &sendgridProvider{apiKey: key}, nil
}

func (p *sendgridProvider) Send(_ context.Context, msg Message) error {
	from := mail.NewEmail(msg.From.Name, msg.From.Email)
	to := mail.NewEmail(msg.To.Name, msg.To.Email)
	content := mail.NewContent("text/html", fallbackContent(msg.HTMLBody))
	m := mail.NewV3MailInit(from, msg.Subject, to, content)
	if msg.TemplateID != "" {
		m.SetTemplateID(msg.TemplateID)
	}
	for k, v := range msg.Data {
		m.Personalizations[0].SetDynamicTemplateData(k, v)
	}
	request := sendgrid.GetRequest(p.apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	resp, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("sendgrid send failed: status=%d body=%s", resp.StatusCode, resp.Body)
	}
	return nil
}

func fallbackContent(body string) string {
	if strings.TrimSpace(body) == "" {
		return "-"
	}
	return body
}
