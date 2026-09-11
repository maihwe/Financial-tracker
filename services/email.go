package services

import (
	"fmt"
	"os"

	"github.com/resend/resend-go/v4"
)

// EmailService sends application emails.
type EmailService struct {
	client *resend.Client
	from   string
}

// NewEmailService creates a new Resend email service.
func NewEmailService() *EmailService {
	apiKey := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("EMAIL_FROM")

	if apiKey == "" {
		return nil
	}

	if from == "" {
		from = "onboarding@resend.dev"
	}

	return &EmailService{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

// SendPasswordResetEmail sends a password-reset link.
func (service *EmailService) SendPasswordResetEmail(
	to string,
	resetURL string,
) error {

	if service == nil {
		return fmt.Errorf("email service is not configured")
	}

	_, err := service.client.Emails.Send(
		&resend.SendEmailRequest{
			From:    service.from,
			To:      []string{to},
			Subject: "Reset your Financial Tracker password",
			Html: fmt.Sprintf(`
				<h2>Password Reset</h2>

				<p>
					We received a request to reset your
					Financial Tracker password.
				</p>

				<p>
					<a href="%s">
						Reset your password
					</a>
				</p>

				<p>
					This link will expire in 30 minutes.
				</p>

				<p>
					If you did not request this, you can
					safely ignore this email.
				</p>
			`, resetURL),
		},
	)

	return err
}
