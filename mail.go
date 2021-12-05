package main

import (
	"log"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Mailer interface {
	SendMail(fromName, fromEmail, toName, toEmail, subject, body string) error
}


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


type logMailer struct{}

func NewLogMailer() *logMailer {
	mailer := logMailer{}
	return &mailer
}

func (m *logMailer) SendMail(fromEmail, fromName, toName, toEmail, subject, body string) error {
	log.Printf("--- LogMailer ---\n")
	log.Printf("fromName:  %s\n", fromName)
	log.Printf("fromEmail: %s\n", fromEmail)
	log.Printf("toName:    %s\n", toName)
	log.Printf("toEmail:   %s\n", toEmail)

	log.Printf("subject: %s\n", subject)
	log.Printf("body:\n  %s\n", strings.Replace(body, "\n", "\n  ", -1))
	return nil
}
