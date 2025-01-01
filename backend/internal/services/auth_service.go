package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/utils"
	"backend/internal/validators"
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

type AuthService struct {
	UserRepo interfaces.UserRepositoryInterface
	AuthRepo interfaces.AuthRepositoryInterface
}

func NewAuthService(userRepo interfaces.UserRepositoryInterface, authRepo interfaces.AuthRepositoryInterface) *AuthService {
	return &AuthService{UserRepo: userRepo, AuthRepo: authRepo}
}

func (s *AuthService) ValidateCredentials(email string, password string) (*models.User, error) {
	if err := validators.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := validators.ValidatePassword(password); err != nil {
		return nil, err
	}

	user, err := s.UserRepo.GetUserByEmail(email)
	if err != nil {
		return nil, validators.ErrUserNotFound
	}

	if !utils.ComparePassword(user.Password, password) {
		return nil, errors.New("invalid credentials")
	}

	user.LastLogin = time.Now()
	if err := s.UserRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", "", fmt.Errorf("SECRET_KEY is not set")
	}

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
		return []byte(secretKey), nil
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
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", "", fmt.Errorf("invalid token structure")
	}

	// Retrieve the user from the repository
	fetchedUser, err := s.UserRepo.GetUserByID(uint(userID))
	if err != nil {
		return "", "", validators.ErrUserNotFound
	}

	// Check if the user is active
	if fetchedUser.IsDeleted {
		return "", "", validators.ErrUserDeleted
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
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return "", "", fmt.Errorf("SECRET_KEY is not set")
	}

	// Access Token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) RevokeRefreshToken(refreshToken string) error {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		return fmt.Errorf("SECRET_KEY is not set")
	}

	// Parse token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return fmt.Errorf("invalid refresh token: %w", err)
	}

	// Validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token claims")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid user_id in token claims")
	}

	// Verify user existence
	user, err := s.UserRepo.GetUserByID(uint(userID))
	if err != nil {
		return validators.ErrUserNotFound
	}
	if user.IsDeleted {
		return validators.ErrUserDeleted
	}

	// Revoke token
	err = s.AuthRepo.RevokeToken(refreshToken)
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
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
