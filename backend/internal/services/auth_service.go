package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"errors"
	"fmt"
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

	user.LastLogin = time.Now()
	if err := s.UserRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token valid for 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key")) // Replace with your actual secret key
}

func (s *AuthService) RefreshToken(refreshToken string) (string, error) {
	// Validate and parse the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil // Use your actual secret key
	})
	if err != nil {
		return "", fmt.Errorf("invalid or expired refresh token")
	}

	// Extract claims and validate the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", fmt.Errorf("invalid token claims")
	}

	// Extract user ID from the claims
	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok {
		return "", fmt.Errorf("invalid token structure")
	}

	// Retrieve the user from the database
	fetchedUser, err := s.UserRepo.GetUserByID(uint(userID))

	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	// Check if the user is active
	if fetchedUser.IsDeleted {
		return "", fmt.Errorf("user is deactivated")
	}

	// Generate a new access token
	newToken, err := s.GenerateToken(fetchedUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate new token: %v", err)
	}

	return newToken, nil
}
