package email

//logic which will send the email 
//using resend for that purpose to send the mail
import (
	"bytes"
	"html/template"
	"github.com/resend/resend-go/v2"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/Nischal07bot/go_boiler_backend/internal/config"
	"github.com/pkg/errors"
)

type Client struct {
	client *resend.Client
	logger *zerolog.Logger
}

func NewClient(cfg *config.Config,logger *zerolog.Logger) *Client {
	return &Client{
		client: resend.NewClient(cfg.Integration.ResendApi),
		logger: logger,
	}
}
func (c *Client) SendEmail(to, subject string, templateName Template, data map[string]string) error {
	tmplPath := fmt.Sprintf("%s/%s", "templates/emails", templateName)

	tmpl,err := template.ParseFiles(tmplPath)
	if err != nil {
		return errors.Wrapf(err, "failed to parse email template: %s", tmplPath)
	}
	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return errors.Wrapf(err, "failed to execute email template: %s", tmplPath)
	}
	params := &resend.SendEmailRequest{
		From: fmt.Sprintf("%s <%s>", "Alfred", "onboarding@resend.dev"),
		To: []string{to},
		Subject: subject,
		Html: body.String(),
	}
	_, err = c.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}