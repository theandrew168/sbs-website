package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type Mailer interface {
	SendMail(fromName, fromEmail, toName, toEmail, subject, body string) error
}

type postmarkMailer struct {
	postmarkAPIKey string
}

func NewPostmarkMailer(postmarkAPIKey string) Mailer {
	m := postmarkMailer{
		postmarkAPIKey: postmarkAPIKey,
	}
	return &m
}

func (m *postmarkMailer) SendMail(fromName, fromEmail, toName, toEmail, subject, body string) error {
	message := struct {
		From     string `json:"From"`
		To       string `json:"To"`
		Subject  string `json:"Subject"`
		TextBody string `json:"TextBody"`
	}{
		From:     fromEmail,
		To:       toEmail,
		Subject:  subject,
		TextBody: body,
	}

	b, err := json.Marshal(message)
	if err != nil {
		return err
	}

	client := &http.Client{}
	endpoint := "https://api.postmarkapp.com/email"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Postmark-Server-Token", m.postmarkAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// error if status isn't a 2xx
	if resp.Status[0] != '2' {
		return fmt.Errorf("failed to send email: %s", resp.Status)
	}

	return nil
}

type logMailer struct{}

func NewLogMailer() Mailer {
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
