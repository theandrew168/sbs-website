package mail

import (
	"log"
	"strings"
)

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
