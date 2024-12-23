package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/utils"
	"backend/internal/validators"
	"errors"
	"os"
)

type UserService struct {
	repo interfaces.UserRepositoryInterface
}

func NewUserService(repo interfaces.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(id uint) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

// User registration (simplified)
func (s *UserService) RegisterUser(user *models.User) error {
	// Validate the user input
	if err := validators.ValidateUser(user); err != nil {
		return err
	}

	// Hash the user's password
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	// Generate verification token
	token := utils.GenerateVerificationToken()
	user.VerificationToken = token
	user.EmailVerified = false

	// Save user to DB
	if err := s.repo.CreateUser(user); err != nil {
		return err
	}

	// Create an email sender
	sender := &models.SmtpEmailSender{}

	if os.Getenv("TESTING") != "true" {
		err := utils.SendVerificationEmail(*user, sender)

		if err != nil {
			return err
		}
	}

	return nil
}

// Email verification
func (s *UserService) VerifyEmail(token string) error {
	user, err := s.repo.GetUserByVerificationToken(token)
	if err != nil {
		return err
	}

	if validators.IsVerificationTokenExpired(user.CreatedAt, user.EmailVerified) {
		return errors.New("invalid token or token expired")
	}

	user.EmailVerified = true
	user.VerificationToken = "" // Clear the token after verification

	// Save changes
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserFields(userID uint, fields interface{}) error {
	// Update logic for general fields like bio, profile picture, etc.
	return s.repo.UpdateUserFields(userID, fields)
}

func (s *UserService) UpdatePassword(userID uint, currentPassword string, newPassword string) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if !utils.ComparePassword(user.Password, currentPassword) {
		return errors.New("current password is incorrect")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.UpdateUser(user)
}

func (s *UserService) UpdateEmail(userID uint, newEmail string) error {
	// Generate a new verification token
	token := utils.GenerateVerificationToken()

	// Update the user with the new email and token
	return s.repo.UpdateUserEmail(userID, newEmail, token)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	return s.repo.GetUserByUsername(username)
}

// User deletion
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.DeleteUser(user); err != nil {
		return err
	}

	return nil
}
