package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"errors"
	"time"
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
func (repo *AuthRepository) RevokeToken(refreshToken string) error {
	var existingToken models.RevokedToken
	err := repo.db.Where("token = ?", refreshToken).First(&existingToken).Error
	if err == nil {
		// Token is already revoked
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected error
		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Occured an error revoking the token", err)
	}

	// Token not found, proceed to insert
	newRevokedToken := models.RevokedToken{
		Token:     refreshToken,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour), // Adjust based on your requirements
	}
	return repo.db.Create(&newRevokedToken).Error
}
