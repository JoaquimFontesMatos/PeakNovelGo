package services

import (
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system environment variables.")
	}
}

var secret = os.Getenv("SECRET_KEY")

type AuthService struct {
	UserRepo repositories.UserRepository
	AuthRepo repositories.AuthRepository
}

func NewAuthService(userRepo repositories.UserRepository, authRepo repositories.AuthRepository) *AuthService {
	return &AuthService{UserRepo: userRepo, AuthRepo: authRepo}
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

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	// Check if the refresh token is in the revoked tokens table
	isRevoked, err := s.AuthRepo.CheckIfTokenRevoked(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("error checking if token is revoked: %v", err)
	}
	if isRevoked {
		return "", "", fmt.Errorf("refresh token has been revoked")
	}

	// Validate and parse the refresh token (same as before)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil // Replace with your actual secret key
	})
	if err != nil {
		return "", "", fmt.Errorf("invalid or expired refresh token")
	}

	// Extract claims and validate the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", fmt.Errorf("invalid token claims")
	}

	// Extract user ID from the claims
	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok {
		return "", "", fmt.Errorf("invalid token structure")
	}

	// Retrieve the user from the repository
	fetchedUser, err := s.AuthRepo.GetUserByID(uint(userID))
	if err != nil {
		return "", "", fmt.Errorf("user not found: %v", err)
	}

	// Check if the user is active
	if fetchedUser.IsDeleted {
		return "", "", fmt.Errorf("user is deactivated")
	}

	// Generate new tokens
	newAccessToken, newRefreshToken, err := s.GenerateToken(fetchedUser)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %v", err)
	}

	// Optionally, revoke the old refresh token
	// You can choose to revoke the refresh token every time a new refresh happens
	err = s.AuthRepo.RevokeToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %v", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) GenerateToken(user *models.User) (string, string, error) {
	// Generate new tokens
	// Access Token (expires in 15 minutes)
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(15 * time.Minute).Unix(), // Access token valid for 15 minutes
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %v", err)
	}

	// Refresh Token (expires in 24 hours)
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(), // Refresh token valid for 24 hours
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte("your_secret_key"))
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %v", err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) RevokeRefreshToken(refreshToken string) error {
	// Validate and parse the refresh token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil // Replace with your actual secret key
	})
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	// Extract claims and validate the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token claims")
	}

	// Extract user ID from the claims (optional, depending on your requirements)
	userID, ok := claims["user_id"].(float64) // JWT stores numbers as float64
	if !ok {
		return fmt.Errorf("invalid token structure")
	}

	fetchedUser, err := s.AuthRepo.GetUserByID(uint(userID))
	if err != nil || fetchedUser.IsDeleted {
		return fmt.Errorf("user is deactivated or not found")
	}

	// Delegate the revocation logic to the repository
	err = s.AuthRepo.RevokeToken(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}

	return nil
}

func (s *AuthService) Logout(refreshToken string) error {
	// Revoke the refresh token
	err := s.RevokeRefreshToken(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}
	return nil
}
