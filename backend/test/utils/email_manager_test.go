package utils_test

import (
	"backend/internal/models"
	"backend/internal/utils"
	"errors"
	"net/smtp"
	"testing"
)

type MockEmailSender struct {
	SendMailCalled bool
	ExpectedError  error
}

func (m *MockEmailSender) SendMail(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// Simulate the behavior of the real SendMail function
	m.SendMailCalled = true
	if m.ExpectedError != nil {
		return m.ExpectedError
	}
	return nil
}

func TestSendVerificationEmail_Success(t *testing.T) {
	// Create a mock email sender
	mockSender := &MockEmailSender{}

	// Sample user data
	user := models.User{
		Username:          "testuser",
		Email:             "testuser@example.com",
		VerificationToken: "12345",
	}

	// Call the function being tested
	err := utils.SendVerificationEmail(user, mockSender)

	// Assert that no error occurred
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	// Assert that SendMail was called
	if !mockSender.SendMailCalled {
		t.Errorf("Expected SendMail to be called, but it wasn't")
	}
}

func TestSendVerificationEmail_Failure(t *testing.T) {
	// Create a mock email sender that simulates an error
	mockSender := &MockEmailSender{
		ExpectedError: errors.New("SMTP server error"),
	}

	// Sample user data
	user := models.User{
		Username:          "testuser",
		Email:             "testuser@example.com",
		VerificationToken: "12345",
	}

	// Call the function being tested
	err := utils.SendVerificationEmail(user, mockSender)

	// Assert that an error occurred
	if err == nil {
		t.Errorf("Expected error, but got none")
	}

	// Assert that SendMail was called
	if !mockSender.SendMailCalled {
		t.Errorf("Expected SendMail to be called, but it wasn't")
	}
}
