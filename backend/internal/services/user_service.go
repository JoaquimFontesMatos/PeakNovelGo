package services

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"backend/internal/validators"
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

// RegisterUser registers a new user.
//
// Parameters:
//   - userFields dtos.RegisterRequest (RegisterRequest struct)
//
// Returns:
//   - CONFLICT_ERROR if the user already exists
//   - INTERNAL_SERVER_ERROR if the user could not be created
func (s *UserService) RegisterUser(userFields *dtos.RegisterRequest) error {
	_, err := s.repo.GetUserByEmail(userFields.Email)
	if err == nil {
		return types.WrapError("CONFLICT_ERROR", "User already registered", nil)
	}

	birthDate, err := time.Parse("2006-01-02", userFields.DateOfBirth)
	if err != nil {
		return types.WrapError("INTERNAL_SERVER_ERROR", "Failed to parse date of birth", err)
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
	sender := &types.SmtpEmailSender{}

	if os.Getenv("TESTING") != "true" {
		err := utils.SendVerificationEmail(user, sender)

		if err != nil {
			return err
		}
	}

	return nil
}

// VerifyEmail verifies the email of a user.
//
// Parameters:
//   - token string (verification token)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user could not be updated
//   - VALIDATION_ERROR if the token is invalid
//   - USER_NOT_FOUND_ERROR if the user is not found
//   - USER_DEACTIVATED_ERROR if the user is deactivated
//   - INVALID_TOKEN_ERROR if the token is invalid
func (s *UserService) VerifyEmail(token string) error {
	if err := validators.ValidateToken(token); err != nil {
		return err
	}

	user, err := s.repo.GetUserByVerificationToken(token)
	if err != nil {
		return types.ErrUserNotFound
	}

	if user.IsDeleted {
		return types.ErrUserDeactivated
	}

	if validators.IsVerificationTokenExpired(user.CreatedAt, user.EmailVerified) {
		return types.ErrTokenExpired
	}

	user.EmailVerified = true
	user.VerificationToken = ""

	// Save changes
	if err := s.repo.UpdateUser(user); err != nil {
		return err
	}

	return nil
}

// UpdateUserFields updates the fields of a user.
//
// Parameters:
//   - userID uint (ID of the user)
//   - fields dtos.UpdateRequest (UpdateRequest struct)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user was not found, deactivated or the update failed
//   - VALIDATION_ERROR if the fields are invalid
//   - USER_NOT_FOUND_ERROR if the user is not found
//   - USER_DEACTIVATED_ERROR if the user is deactivated
func (s *UserService) UpdateUserFields(userID uint, fields dtos.UpdateRequest) error {
	if err := validators.ValidateUserFields(fields); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return types.ErrUserNotFound
	}

	if user.IsDeleted {
		return types.ErrUserDeactivated
	}

	return s.repo.UpdateUserFields(userID, fields)
}

// UpdatePassword updates the password of a user.
//
// Parameters:
//   - userID uint (ID of the user)
//   - currentPassword string (current password of the user)
//   - newPassword string (new password of the user)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user was not found, deactivated or the update failed
//   - VALIDATION_ERROR if the new password is not valid
//   - USER_NOT_FOUND_ERROR if the user is not found
//   - USER_DEACTIVATED_ERROR if the user is deactivated
//   - PASSWORD_DIFF_ERROR if the new password is the same as the current password
//   - INVALID_PASSWORD_ERROR if the current password is invalid
func (s *UserService) UpdatePassword(userID uint, currentPassword string, newPassword string) error {
	if err := validators.ValidatePassword(newPassword); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return types.ErrUserNotFound
	}

	if user.IsDeleted {
		return types.ErrUserDeactivated
	}

	if !utils.ComparePassword(user.Password, currentPassword) {
		return types.ErrInvalidPassword
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

// UpdateEmail updates the email of a user.
//
// Parameters:
//   - userID uint (ID of the user)
//   - newEmail string (new email of the user)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if the user was not found, deactivated or the update failed
//   - VALIDATION_ERROR if the new email is not valid
//   - USER_NOT_FOUND_ERROR if the user is not found
//   - USER_DEACTIVATED_ERROR if the user is deactivated
func (s *UserService) UpdateEmail(userID uint, newEmail string) error {
	if err := validators.ValidateEmail(newEmail); err != nil {
		return err
	}

	// Check if user exists
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return types.ErrUserNotFound
	}

	if user.IsDeleted {
		return types.ErrUserDeactivated
	}

	// Generate a new verification token
	token := utils.GenerateVerificationToken()

	// Update the user with the new email and token
	return s.repo.UpdateUserEmail(userID, newEmail, token)
}

// GetUserByEmail gets a user by email.
//
// Parameters:
//   - email string (email of the user)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	if err := validators.ValidateEmail(email); err != nil {
		return nil, err
	}

	return s.repo.GetUserByEmail(email)
}

// GetUserByUsername gets a user by username.
//
// Parameters:
//   - username string (username of the user)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
func (s *UserService) GetUserByUsername(username string) (*models.User, error) {
	if err := validators.ValidateUsername(username); err != nil {
		return nil, err
	}

	return s.repo.GetUserByUsername(username)
}

// DeleteUser deactivates a user in the database.
//
// Parameters:
//   - id uint (ID of the user to deactivate)
//
// Returns:
//   - USER_NOT_FOUND_ERROR if the user is not found
//   - USER_DEACTIVATED_ERROR if the user is deactivated
//   - INTERNAL_SERVER_ERROR if the user could not be deactivated
func (s *UserService) DeleteUser(id uint) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return types.ErrUserNotFound
	}

	if err := validators.ValidateIsDeleted(*user); err != nil {
		return err
	}

	if err := s.repo.DeleteUser(user); err != nil {
		return err
	}

	return nil
}
