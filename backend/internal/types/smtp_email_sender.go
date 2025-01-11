package types

import (
	"net/smtp"
)

// SmtpEmailSender struct represents an SMTP email sender.
type SmtpEmailSender struct{}

// SendMail sends an email using SMTP.
//
// Parameters:
//   - addr string (SMTP server address)
//   - auth smtp.Auth (SMTP authentication)
//   - from string (email address of the sender)
//   - to []string (list of email addresses to send the email to)
//   - msg []byte (email message)
//
// Returns:
//   - EMAIL_SEND_ERROR if there is an error sending the email
func (sender *SmtpEmailSender) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Send the email using SMTP
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return WrapError(EMAIL_SEND_ERROR, "Failed to send email", err)
	}
	return nil
}
