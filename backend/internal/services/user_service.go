package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"backend/internal/validators"
	"errors"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.Repo.GetUserByID(id)
}

// User registration (simplified)
func (s *UserService) RegisterUser(user *models.User) error {
	// Generate verification token
	token := utils.GenerateVerificationToken()
	user.VerificationToken = token
	user.EmailVerified = false

	// Save user to DB
	if err := s.Repo.CreateUser(user); err != nil {
		return err
	}

	// Create an email sender
	sender := &models.SmtpEmailSender{}

	// Send verification email (you would implement the actual email sending)
	err := utils.SendVerificationEmail(*user, sender)

	if err != nil {
		return err
	}

	return nil
}

// Email verification
func (s *UserService) VerifyEmail(token string) error {
	user, err := s.Repo.GetUserByVerificationToken(token)
	if err != nil {
		return err
	}

	if validators.IsVerificationTokenExpired(user.CreatedAt, user.EmailVerified) {
		return errors.New("invalid token or token expired")
	}

	user.EmailVerified = true
	user.VerificationToken = "" // Clear the token after verification

	// Save changes
	if err := s.Repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	// Business logic before saving
	return s.Repo.UpdateUser(user)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.Repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.Repo.GetUserByUsername(username)
}

// User deletion
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.Repo.DeleteUser(user); err != nil {
		return err
	}

	return nil
}
