package email

import (
	"fmt"
	"os"
)

type EmailService struct {
	// In production, you would use an email provider like SendGrid, AWS SES, etc.
	// For now, we'll just log the emails
	fromEmail string
	apiKey    string
}

func NewEmailService() *EmailService {
	return &EmailService{
		fromEmail: os.Getenv("EMAIL_FROM"),
		apiKey:    os.Getenv("EMAIL_API_KEY"),
	}
}

func (s *EmailService) SendVerificationEmail(toEmail, name, token string) error {
	verifyURL := fmt.Sprintf("%s/verify-email?token=%s", os.Getenv("FRONTEND_URL"), token)
	
	subject := "Verify your email address"
	body := fmt.Sprintf(`
Hello %s,

Thank you for signing up! Please verify your email address by clicking the link below:

%s

This link will expire in 24 hours.

If you didn't create an account, please ignore this email.

Best regards,
API Platform Team
`, name, verifyURL)

	// TODO: Implement actual email sending using your preferred provider
	// For now, just log it
	fmt.Printf("=== EMAIL ===\nTo: %s\nSubject: %s\nBody:%s\n=============\n", toEmail, subject, body)
	
	return nil
}

func (s *EmailService) SendPasswordResetEmail(toEmail, name, token string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", os.Getenv("FRONTEND_URL"), token)
	
	subject := "Reset your password"
	body := fmt.Sprintf(`
Hello %s,

We received a request to reset your password. Click the link below to reset it:

%s

This link will expire in 1 hour.

If you didn't request a password reset, please ignore this email.

Best regards,
API Platform Team
`, name, resetURL)

	// TODO: Implement actual email sending
	fmt.Printf("=== EMAIL ===\nTo: %s\nSubject: %s\nBody:%s\n=============\n", toEmail, subject, body)
	
	return nil
}

func (s *EmailService) SendWelcomeEmail(toEmail, name string) error {
	subject := "Welcome to API Platform!"
	body := fmt.Sprintf(`
Hello %s,

Welcome to API Platform! Your email has been verified and your account is now active.

You can now:
- Create and deploy REST APIs
- Browse the API marketplace
- Monitor your API usage and analytics

Get started by logging in at: %s

Best regards,
API Platform Team
`, name, os.Getenv("FRONTEND_URL"))

	// TODO: Implement actual email sending
	fmt.Printf("=== EMAIL ===\nTo: %s\nSubject: %s\nBody:%s\n=============\n", toEmail, subject, body)
	
	return nil
}
