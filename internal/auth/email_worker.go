package auth

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"
	"todoApp3/internal/domain"
	"todoApp3/internal/email"
)

func StartEmailWorker(ctx context.Context, ch <-chan domain.EmailJob, wg *sync.WaitGroup, logger *slog.Logger) {
	defer wg.Done()

	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <style>
        body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
        .container { max-width: 600px; margin: 20px auto; background: #ffffff; border-radius: 8px; box-shadow: 0 0 10px rgba(0,0,0,0.1); overflow: hidden; }
        .header { background-color: #4A90E2; color: #ffffff; padding: 20px; text-align: center; }
        .content { padding: 30px 20px; color: #333333; line-height: 1.6; }
        .message-box { background-color: #f9f9f9; border-left: 4px solid #4A90E2; padding: 15px; margin: 20px 0; font-style: italic; }
        .footer { background-color: #eeeeee; padding: 15px; text-align: center; font-size: 12px; color: #777777; }
        .button { display: inline-block; padding: 10px 20px; margin-top: 20px; background-color: #4A90E2; color: white; text-decoration: none; border-radius: 5px; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Welcome to TodoApp! 🚀</h1>
        </div>

        <div class="content">
            <p>Dear <strong>%s %s</strong>,</p>
            
            <p>We are thrilled to have you on board! Your registration has been completed successfully.</p>
            
            <div class="message-box">
                "%s"
            </div>

            <p>You can click the button below to start using your account immediately.</p>
            
            <center>
                <a href="#" class="button">Go to App</a>
            </center>
        </div>

        <div class="footer">
            <p>© 2024 TodoApp. All rights reserved.</p>
            <p>This is an automated message, please do not reply.</p>
        </div>
    </div>
</body>
</html>
`
	for job := range ch {

		body := fmt.Sprintf(htmlTemplate, job.Name, job.Surname, job.Message)

		smtpReq := email.SmtpRequest{
			To:      job.Email,
			Subject: "Welcome",
			Body:    body,
		}

		jobCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		err := email.Send(jobCtx, &smtpReq, logger)
		cancel()

		if err != nil {
			logger.Error("email send failed", "to", job.Email, "err", err)
			continue
		}

	}

	logger.Info("email worker channel closed, exiting")
}
