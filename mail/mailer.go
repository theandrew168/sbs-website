package mail

type Mailer interface {
	SendMail(fromName, fromEmail, toName, toEmail, subject, body string) error
}
