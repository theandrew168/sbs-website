package mail

import (
	"log"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type sendGridMailer struct {
	sendGridAPIKey string
}

func NewSendGridMailer(sendGridAPIKey string) *sendGridMailer {
	mailer := sendGridMailer{
		sendGridAPIKey: sendGridAPIKey,
	}
	return &mailer
}

func (m *sendGridMailer) SendMail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	from := mail.NewEmail(fromName, fromEmail)
	to := mail.NewEmail(toName, toEmail)

	message := mail.NewSingleEmailPlainText(from, subject, to, body)

	client := sendgrid.NewSendClient(m.sendGridAPIKey)
	response, err := client.Send(message)
	if err != nil {
		return err
	}

	log.Printf("SendGridMailer: message sent to %s, status code %d\n", toEmail, response.StatusCode)
	return nil
}
