package email

import (
	"context"
	"log/slog"
	"time"

	"github.com/wneessen/go-mail"
)

type SmtpRequest struct {
	To      string
	Subject string
	Body    string
}

func Send(ctx context.Context, req *SmtpRequest, logger *slog.Logger) error {
	m := mail.NewMsg()
	if err := m.From("berkay_test@test.com"); err != nil {
		logger.Error("failed to set from address",
			slog.String("message", err.Error()))
		return err
	}
	if err := m.To(req.To); err != nil {
		logger.Error("failed to set to address",
			slog.String("message", err.Error()))
		return err
	}
	m.Subject(req.Subject)
	m.SetBodyString(mail.TypeTextHTML, req.Body)

	c, err := mail.NewClient("sandbox.smtp.mailtrap.io",
		mail.WithPort(2525),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername("40247f05c4d5be"),
		mail.WithPassword("9e27e49a87793e"),
		mail.WithTimeout(10*time.Second),
	)
	if err != nil {
		logger.Error("failed to create mail client",
			slog.String("message", err.Error()))
		return err
	}

	if err := c.DialAndSendWithContext(ctx, m); err != nil {
		logger.Error("failed to send email",
			slog.String("message", err.Error()))
		return err
	}
	return nil

}
