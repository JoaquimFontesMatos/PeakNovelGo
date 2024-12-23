package interfaces

import "backend/internal/models"

type UserServiceInterface interface {
	GetUser(id uint) (*models.User, error)
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	DeleteUser(id uint) error
	VerifyEmail(token string) error
	UpdateUserFields(userID uint, fields models.UpdateFields) error
	UpdatePassword(userID uint, currentPassword string, newPassword string) error
	UpdateEmail(userID uint, newEmail string) error
}
