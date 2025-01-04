package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/utils"
	"backend/internal/types"

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
	UserRepo  interfaces.UserRepositoryInterface
	AuthRepo  interfaces.AuthRepositoryInterface
	SecretKey []byte
}

func NewAuthService(userRepo interfaces.UserRepositoryInterface, authRepo interfaces.AuthRepositoryInterface) *AuthService {
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY environment variable not set")
	}
	return &AuthService{
		UserRepo:  userRepo,
		AuthRepo:  authRepo,
		SecretKey: []byte(secretKey),
	}
}

func (s *AuthService) ValidateCredentials(email string, password string) (*models.User, error) {
	if err := validators.ValidateEmail(email); err != nil {
		return nil, err
	}
	if err := validators.ValidatePassword(password); err != nil {
		return nil, err
	}

	userChan := make(chan *models.User)
	errChan := make(chan error)
	defer close(userChan)
	defer close(errChan)

	go func() {
		user, err := s.UserRepo.GetUserByEmail(email)
		if err != nil {
			errChan <- types.ErrUserNotFound
			return
		}
		userChan <- user
	}()

	select {
	case user := <-userChan:
		if !utils.ComparePassword(user.Password, password) {
			return nil, errors.New("invalid credentials")
		}
		user.LastLogin = time.Now()
		if err := s.UserRepo.UpdateUser(user); err != nil {
			return nil, err
		}
		return user, nil
	case err := <-errChan:
		return nil, err
	}
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
		return s.SecretKey, nil
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
		return "", "", types.ErrUserNotFound
	}

	// Check if the user is active
	if fetchedUser.IsDeleted {
		return "", "", types.ErrUserDeleted
	}

	revokeErrChan := make(chan error)
	tokenChan := make(chan [2]string)
	defer close(revokeErrChan)
	defer close(tokenChan)

	go func() {
		revokeErrChan <- s.AuthRepo.RevokeToken(refreshToken)
	}()

	go func() {
		accessToken, refreshToken, err := s.GenerateToken(fetchedUser)
		if err != nil {
			tokenChan <- [2]string{"", ""}
			revokeErrChan <- err
			return
		}
		tokenChan <- [2]string{accessToken, refreshToken}
	}()

	revokeErr := <-revokeErrChan
	if revokeErr != nil {
		return "", "", fmt.Errorf("failed to revoke token: %v", revokeErr)
	}

	tokens := <-tokenChan
	if tokens[0] == "" && tokens[1] == "" {
		return "", "", fmt.Errorf("failed to generate tokens")
	}

	return tokens[0], tokens[1], nil
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
		return types.ErrUserNotFound
	}
	if user.IsDeleted {
		return types.ErrUserDeleted
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
