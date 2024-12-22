package utils

import (
	"fmt"
	"log"
	"net/smtp"
	"os"

	"backend/internal/models"

	"github.com/joho/godotenv"
)

// EmailSender interface allows mocking of SMTP SendMail
type EmailSender interface {
	SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

// SendVerificationEmail sends a verification email to the user
func SendVerificationEmail(user models.User, sender EmailSender) error {
	// Create the verification URL
	verificationURL := fmt.Sprintf("http://your-app.com/verify-email?token=%s", user.VerificationToken)

	// Compose the email content
	subject := "Email Verification"
	body := fmt.Sprintf(`
		<p>Hello %s,</p>
		<p>Thank you for registering with us. Please verify your email address by clicking the link below:</p>
		<a href="%s">%s</a>
		<p>If you did not request this, please ignore this email.</p>
	`, user.Username, verificationURL, verificationURL)

	// Set up the email message
	from := os.Getenv("SMTP_USERNAME")
	to := []string{user.Email}
	message := []byte("To: " + user.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/html; charset=UTF-8\r\n" +
		"\r\n" +
		body + "\r\n")

	// Create SMTP auth
	auth := smtp.PlainAuth("", from, os.Getenv("SMTP_PASSWORD"), os.Getenv("SMTP_HOST"))

	// Check if sender is provided, if not, fall back to default sender
	if sender == nil {
		// Fall back to a default sender, like an SMTP sender
		sender = &models.SmtpEmailSender{}
	}

	// Send the email using the provided EmailSender
	err := sender.SendMail(os.Getenv("SMTP_HOST")+":"+os.Getenv("SMTP_PORT"), auth, from, to, message)
	if err != nil {
		log.Printf("Failed to send verification email: %v", err)
		return err
	}

	log.Println("Verification email sent successfully to", user.Email)
	return nil
}
