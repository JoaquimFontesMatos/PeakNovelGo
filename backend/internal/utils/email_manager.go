package utils

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"

	"backend/internal/interfaces"
	"backend/internal/models"
	"backend/internal/types"
	"backend/internal/types/errors"
)

// SendVerificationEmail sends a verification email to the given user.
//
// It generates a verification URL based on the FRONTEND_URL environment variable (or defaults to "http://localhost:3000") and the user's verification token.
// The email includes a link to the verification URL, which the user can click to activate their account.
// If an EmailSender is provided, it uses that to send the email. Otherwise, it falls back to a default SMTP sender using environment variables for SMTP configuration.
//
// Parameters:
//  - user (models.User) - The user to whom the verification email should be sent.
//  - sender (interfaces.EmailSender) - The EmailSender interface to use for sending the email. If nil, a default SMTP sender will be used.
//
// Returns:
//  - error - An error object indicating whether the email was sent successfully.  Returns types.WrapError(errors.EMAIL_SEND, "Failed to send verification email", http.StatusInternalServerError, err) if the email fails to send.
//
// Generated on 2024-07-26
func SendVerificationEmail(user models.User, sender interfaces.EmailSender) error {
	baseUrl := os.Getenv("FRONTEND_URL")
	if baseUrl == "" {
		baseUrl = "http://localhost:3000"
	}

	verificationURL := fmt.Sprintf("%s/auth/activate-account/%s", baseUrl, user.VerificationToken)

	subject := "Email Verification"
	escapedUsername := template.HTMLEscapeString(user.Username)
	escapedVerificationURL := template.HTMLEscapeString(verificationURL)
	body := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thank you for registering with us. Please verify your email address by clicking the link below:</p>
		<a href="%s">%s</a>
		<p>If you did not request this, please ignore this email.</p>
	`, escapedUsername, escapedVerificationURL, escapedVerificationURL)

	from := os.Getenv("SMTP_USERNAME")
	to := []string{user.Email}
	message := []byte("To: " + user.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	if sender == nil {
		sender = &types.SmtpEmailSender{}
	}

	err := sender.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send verification email: %v", err)
		return types.WrapError(errors.EMAIL_SEND, "Failed to send verification email", http.StatusInternalServerError, err)
	}

	log.Println("Verification email sent successfully to", user.Email)
	return nil
}