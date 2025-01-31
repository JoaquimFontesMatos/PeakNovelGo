package interfaces

import "backend/internal/models"
import "backend/internal/dtos"

type UserServiceInterface interface {
	GetUser(id uint) (*models.User, error)
	RegisterUser(user *dtos.RegisterRequest) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	DeleteUser(id uint) error
	VerifyEmail(token string) error
	UpdateUserFields(userID uint, fields dtos.UpdateRequest) error
	UpdatePassword(userID uint, currentPassword string, newPassword string) error
	UpdateEmail(userID uint, newEmail string) error
}
