package repositories

import (
	"backend/internal/models"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CheckIfTokenRevoked checks if the given refresh token has been revoked.
func (r *AuthRepository) CheckIfTokenRevoked(refreshToken string) (bool, error) {
	var revokedToken models.RevokedToken
	err := r.db.Where("token = ?", refreshToken).First(&revokedToken).Error
	if err == nil {
		// Token is revoked
		return true, nil
	} else if err != gorm.ErrRecordNotFound {
		// If some unexpected error occurs while querying the database
		return false, fmt.Errorf("failed to check revocation status: %v", err)
	}
	return false, nil
}

// RevokeToken stores the revoked refresh token in the database.
func (r *AuthRepository) RevokeToken(refreshToken string) error {
	// Parse the refresh token to get its expiration time
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return fmt.Errorf("invalid refresh token")
	}

	// Extract the claims to get the expiration time
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return fmt.Errorf("invalid token claims")
	}

	// Extract expiration time and create a RevokedToken record
	expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)

	revokedToken := models.RevokedToken{
		Token:     refreshToken,
		ExpiredAt: expirationTime,
	}

	// Store the revoked token in the database
	if err := r.db.Create(&revokedToken).Error; err != nil {
		return fmt.Errorf("failed to revoke token: %v", err)
	}

	return nil
}