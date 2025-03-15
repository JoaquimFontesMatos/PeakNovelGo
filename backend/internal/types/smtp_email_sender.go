package types

import (
	"net/smtp"
)

// SmtpEmailSender struct represents an SMTP email sender.
type SmtpEmailSender struct{}

// SendMail sends an email using the provided SMTP configuration.
//
// Parameters:
//   - addr (string): The address of the SMTP server.
//   - auth (smtp.Auth): The authentication mechanism for the SMTP server.
//   - from (string): The sender's email address.
//   - to ([]string): A slice of recipient email addresses.
//   - msg ([]byte): The email message as a byte slice.
//
// Returns:
//   - error: An error if the email could not be sent, nil otherwise.
func (sender *SmtpEmailSender) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}
