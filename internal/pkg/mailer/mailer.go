package mailer

import (
	"github.com/artfoxe6/quick-gin/internal/app/config"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type templateItem struct {
	Id      string
	Subject string
}

var Template = map[string]templateItem{
	"code": {
		Id:      "xxxx-xxxxx",
		Subject: "Verification code",
	},
}

type Email struct {
	templateId string
	subject    string
	from       *mail.Email
	data       map[string]any
}

func New(template templateItem, data map[string]any) Email {
	return Email{
		templateId: template.Id,
		subject:    template.Subject,
		from:       mail.NewEmail("test", "test@gmail.com"),
		data:       data,
	}
}

func (e Email) SendTo(name, email string) (*rest.Response, error) {
	to := mail.NewEmail(name, email)
	content := mail.NewContent("text/html", "-")
	m := mail.NewV3MailInit(e.from, e.subject, to, content)
	m.SetTemplateID(e.templateId)
	for k, v := range e.data {
		m.Personalizations[0].SetDynamicTemplateData(k, v)
	}
	request := sendgrid.GetRequest(config.Sendgrid.Key, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	request.Body = mail.GetRequestBody(m)
	return sendgrid.API(request)
}
