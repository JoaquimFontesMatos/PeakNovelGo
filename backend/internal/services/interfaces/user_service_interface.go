package interfaces

import "backend/internal/models"

type UserServiceInterface interface {
	GetUser(id uint) (*models.User, error)
	RegisterUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	VerifyEmail(token string) error
}
