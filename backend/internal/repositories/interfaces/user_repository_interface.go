package interfaces

import "backend/internal/models"

type UserRepositoryInterface interface {
	CreateUser(user *models.User) error
	GetUserByID(id uint) (*models.User, error)
	GetUserByVerificationToken(token string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByUsername(username string) (*models.User, error)
	UpdateUserEmail(userID uint, newEmail string, token string) error
	UpdateUserFields(userID uint, fields interface{}) error
}
