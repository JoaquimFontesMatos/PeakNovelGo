package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

// AuthRepository struct represents a repository for auth management.
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new AuthRepository instance
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *AuthRepository (pointer to AuthRepository struct)
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CheckIfTokenRevoked checks if the the refresh token is revoked
//
// Parameters:
//   - refreshToken string (refresh token to check)
//
// Returns:
//   - bool (true if the token is revoked, false otherwise)
func (r *AuthRepository) CheckIfTokenRevoked(refreshToken string) bool {
	var revokedToken models.RevokedToken
	err := r.db.Where("token = ?", refreshToken).First(&revokedToken).Error
	return err == nil
}

// RevokeToken revokes the given refresh token so it can't be used again
//
// Parameters:
//   - refreshToken string (refresh token to be revoked)
//
// Returns:
//   - INTERNAL_SERVER_ERROR if an error occurred while revoking token
//   - VALIDATION_ERROR if an the token is invalid
func (r *AuthRepository) RevokeToken(refreshToken string) error {
	// Parse the refresh token to get its expiration time
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, types.WrapError("VALIDATION_ERROR", "Invalid token", fmt.Errorf("unexpected signing method: %v", token.Header["alg"]))
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return types.WrapError("VALIDATION_ERROR", "Invalid token", err)
	}

	// Extract the claims to get the expiration time
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return types.WrapError("VALIDATION_ERROR", "Invalid token", err)
	}

	// Extract expiration time and create a RevokedToken record
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	revokedToken := models.RevokedToken{
		Token:     refreshToken,
		ExpiredAt: expirationTime,
	}

	// Store the revoked token in the database
	if err := r.db.Create(&revokedToken).Error; err != nil {
		return types.WrapError("INTERNAL_SERVER_ERROR", "Failed to revoke token", err)
	}

	return nil
}
