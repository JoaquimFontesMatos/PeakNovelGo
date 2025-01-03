package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/utils"
	"backend/internal/validators"
	"errors"
	"os"
	"time"
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
func (s *UserService) RegisterUser(userFields *models.RegisterRequest) error {
	_, err := s.repo.GetUserByEmail(userFields.Email)
	if err == nil {
		return &validators.ValidationError{Message: "user already registered"}
	}

	birthDate, err := time.Parse("2006-01-02", userFields.DateOfBirth)
	if err != nil {
		return err
	}

	user := models.User{
		Username:       userFields.Username,
		Email:          userFields.Email,
		Password:       userFields.Password,
		Bio:            userFields.Bio,
		ProfilePicture: userFields.ProfilePicture,
		DateOfBirth:    birthDate,
	}

	// Validate the user input
	if err := validators.ValidateUser(&user); err != nil {
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
	if err := s.repo.CreateUser(&user); err != nil {
		return err
	}

	// Create an email sender
	sender := &models.SmtpEmailSender{}

	if os.Getenv("TESTING") != "true" {
		err := utils.SendVerificationEmail(user, sender)

		if err != nil {
			return err
		}
	}

	return nil
}

// Email verification
func (s *UserService) VerifyEmail(token string) error {
	if err := validators.ValidateToken(token); err != nil {
		return err
	}

	user, err := s.repo.GetUserByVerificationToken(token)
	if err != nil {
		return validators.ErrUserNotFound
	}

	if user.IsDeleted {
		return validators.ErrUserDeleted
	}

	if validators.IsVerificationTokenExpired(user.CreatedAt, user.EmailVerified) {
		return validators.ErrTokenExpired
	}

	user.EmailVerified = true
	user.VerificationToken = ""

	// Save changes
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserFields(userID uint, fields models.UpdateFields) error {
	if err := validators.ValidateUserFields(fields); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return validators.ErrUserNotFound
	}

	if user.IsDeleted {
		return validators.ErrUserDeleted
	}

	return s.repo.UpdateUserFields(userID, fields)
}

func (s *UserService) UpdatePassword(userID uint, currentPassword string, newPassword string) error {
	if err := validators.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return validators.ErrUserNotFound
	}

	if user.IsDeleted {
		return validators.ErrUserDeleted
	}

	if !utils.ComparePassword(user.Password, currentPassword) {
		return validators.ErrInvalidPassword
	}

	if err := validators.ValidateIsNewPasswordTheSame(currentPassword, newPassword); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.UpdateUser(user)
}

func (s *UserService) UpdateEmail(userID uint, newEmail string) error {
	if err := validators.ValidateEmail(newEmail); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return validators.ErrUserNotFound
	}

	if user.IsDeleted {
		return validators.ErrUserDeleted
	}

	// Generate a new verification token
	token := utils.GenerateVerificationToken()

	// Update the user with the new email and token
	return s.repo.UpdateUserEmail(userID, newEmail, token)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if err := validators.ValidateEmail(email); err != nil {
		return nil, err
	}

	return s.repo.GetUserByEmail(email)
}

func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	if err := validators.ValidateUsername(username); err != nil {
		return nil, err
	}

	return s.repo.GetUserByUsername(username)
}

// User deletion
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return err
	}

	if err := s.repo.DeleteUser(user); err != nil {
		return err
	}

	return nil
}
