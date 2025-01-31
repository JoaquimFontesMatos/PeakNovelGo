package interfaces

import "backend/internal/models"

type AuthServiceInterface interface {
	ValidateCredentials(email string, password string) (*models.User, error)
	GenerateToken(user *models.User) (string, string, error)
	RefreshToken(refreshToken string) (string, string, *models.User, error)
	Logout(refreshToken string) error
}
