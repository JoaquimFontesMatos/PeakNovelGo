package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type AuthService struct {
	UserRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (s *AuthService) ValidateCredentials(email, password string) (*models.User, error) {
	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil || !utils.ComparePassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) /*TODO: Change this to your secret key*/
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Implement refresh logic
	return "", nil
}
