package email

import (
	"gopkg.in/gomail.v2"
)

type SmtpRequest struct {
	To      string
	Subject string
	Body    string
}

func Send(req *SmtpRequest) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "berkay@test.com")
	m.SetHeader("To", req.To)
	m.SetHeader("Subject", req.Subject)
	m.SetBody("text/html", req.Body)

	d := gomail.NewDialer(
		"sandbox.smtp.mailtrap.io",
		2525,
		"40247f05c4d5be",
		"9e27e49a87793e",
	)

	return d.DialAndSend(m)
}
