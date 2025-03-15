package interfaces

import "net/smtp"

// EmailSender is an interface that defines a method for sending emails.
type EmailSender interface {
	// SendMail sends an email using the provided parameters.
	//
	// Parameters:
	//   - addr string (The address of the SMTP server)
	//   - auth smtp.Auth (The authentication method to use)
	//   - from string (The email address of the sender).
	//   - to []string (The email addresses of the recipients)
	//   - msg []byte (The message to send)
	//
	// Returns:
	//   - error (An error if the email could not be sent)
	SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error
}
