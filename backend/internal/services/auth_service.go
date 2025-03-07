package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/types/errors"
	"backend/internal/utils"
	"net/http"

	"backend/internal/validators"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// AuthService struct represents an authentication service.
type AuthService struct {
	UserRepo  interfaces.UserRepositoryInterface
	AuthRepo  interfaces.AuthRepositoryInterface
	SecretKey []byte
}

// NewAuthService creates a new AuthService instance.
//
// Parameters:
//   - userRepo interfaces.UserRepositoryInterface (UserRepository instance)
//   - authRepo interfaces.AuthRepositoryInterface (AuthRepository instance)
//
// Returns:
//   - *AuthService (pointer to AuthService struct)
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

// ValidateCredentials validates the credentials and returns the user if the credentials are valid.
//
// Parameters:
//   - email string (email of the user)
//   - password string (password of the user)
//
// Returns:
//   - *models.User (pointer to User struct)
//   - INVALID_CREDENTIALS_ERROR if the credentials are invalid
//   - INTERNAL_SERVER_ERROR if the user could not be fetched
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
			errChan <- types.WrapError(errors.INVALID_CREDENTIALS, "Invalid credentials", http.StatusUnauthorized, err)
			return
		}
		userChan <- user
	}()

	select {
	case user := <-userChan:
		if !utils.ComparePassword(user.Password, password) {
			return nil, errors.ErrInvalidCredentials
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

// RefreshToken refreshes the access token and refresh token.
//
// Parameters:
//   - refreshToken string (refresh token to be refreshed)
//
// Returns:
//   - [2]string ([0] is the access token, [1] is the refresh token)
//   - REFRESH_TOKEN_REVOKED_ERROR if the refresh token has been revoked
//   - INTERNAL_SERVER_ERROR if an error occurred while refreshing the token
func (s *AuthService) RefreshToken(refreshToken string) (string, string, *models.User, error) {
	// Check if the refresh token is in the revoked tokens table
	isRevoked := s.AuthRepo.CheckIfTokenRevoked(refreshToken)
	if isRevoked {
		return "", "", nil, errors.ErrRefreshTokenRevoked
	}

	// Validate and parse the refresh token (same as before)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, types.WrapError(errors.INVALID_TOKEN, "Unexpected signing method", http.StatusUnauthorized, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return s.SecretKey, nil
	})
	if err != nil {
		return "", "", nil, types.WrapError(errors.INVALID_TOKEN, "Invalid or expired refresh token", http.StatusUnauthorized, err)
	}

	// Extract claims and validate the token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", nil, types.WrapError(errors.INVALID_TOKEN, "Invalid token claims", http.StatusUnauthorized, err)
	}

	// Extract user ID from the claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return "", "", nil, types.WrapError(errors.INVALID_TOKEN, "Invalid token structure", http.StatusUnauthorized, err)
	}

	// Retrieve the user from the repository
	fetchedUser, err := s.UserRepo.GetUserByID(uint(userID))
	if err != nil {
		return "", "", nil, err
	}

	// Check if the user is active
	if fetchedUser.IsDeleted {
		return "", "", nil, err
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
		return "", "", nil, types.WrapError(errors.REVOKING_TOKEN, "Failed to revoke token", http.StatusInternalServerError, err)
	}

	tokens := <-tokenChan
	if tokens[0] == "" && tokens[1] == "" {
		return "", "", nil, types.WrapError(errors.GENERATING_TOKEN, "Failed to generate tokens", http.StatusInternalServerError, err)
	}

	return tokens[0], tokens[1], fetchedUser, nil
}

// GenerateToken generates the access token and refresh token.
//
// Parameters:
//   - user *models.User (User struct)
//
// Returns:
//   - [2]string ([0] is the access token, [1] is the refresh token)
//   - INTERNAL_SERVER_ERROR if an error occurred while generating the tokens
func (s *AuthService) GenerateToken(user *models.User) (string, string, error) {
	// Access Token
	accessClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.SecretKey)
	if err != nil {
		return "", "", types.WrapError(errors.GENERATING_TOKEN, "Failed to sign access token", http.StatusInternalServerError, err)
	}

	// Refresh Token
	refreshClaims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"roles":   user.Roles,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString(s.SecretKey)
	if err != nil {
		return "", "", types.WrapError(errors.GENERATING_TOKEN, "Failed to sign refresh token",  http.StatusInternalServerError, err)
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) RevokeRefreshToken(refreshToken string) error {
	// Parse token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, types.WrapError(errors.INVALID_TOKEN, "Unexpected signing method", http.StatusUnauthorized, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return s.SecretKey, nil
	})
	if err != nil {
		return types.WrapError(errors.INVALID_TOKEN, "Invalid refresh token", http.StatusUnauthorized, err)
	}

	// Validate claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return types.WrapError(errors.INVALID_TOKEN, "Invalid token claims", http.StatusUnauthorized, err)
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return types.WrapError(errors.INVALID_TOKEN, "Invalid user_id in token claims", http.StatusUnauthorized, err)
	}

	// Verify user existence
	user, err := s.UserRepo.GetUserByID(uint(userID))
	if err != nil {
		return err
	}
	if user.IsDeleted {
		return err
	}

	// Revoke token
	err = s.AuthRepo.RevokeToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}

// Logout revokes the refresh token.
//
// Parameters:
//   - refreshToken string (refresh token to be revoked)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if an error occurred while revoking the token
func (s *AuthService) Logout(refreshToken string) error {
	// Revoke the refresh token
	err := s.RevokeRefreshToken(refreshToken)
	if err != nil {
		return err
	}
	return nil
}
