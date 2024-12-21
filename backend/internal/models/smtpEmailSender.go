package models

import "net/smtp"

type SmtpEmailSender struct{}

// Implement the SendMail method for SmtpEmailSender to fulfill the EmailSender interface
func (sender *SmtpEmailSender) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Send the email using SMTP
	err := smtp.SendMail(addr, auth, from, to, msg)
	if err != nil {
		return err
	}
	return nil
}
