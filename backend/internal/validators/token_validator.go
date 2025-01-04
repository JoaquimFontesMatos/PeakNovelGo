package validators

import (
	"backend/internal/types"

	"time"
)

func IsVerificationTokenExpired(createdAt time.Time, emailVerified bool) bool {
	timeDifference := time.Since(createdAt)

	// Check if the token is expired
	if timeDifference >= 60*time.Minute || emailVerified {
		return true
	}
	return false
}

func ValidateToken(token string) error {
	if token == "" {
		return &types.ValidationError{Message: "token is required"}
	}
	if len(token) > 255 {
		return &types.ValidationError{Message: "token cannot be longer than 255 characters"}
	}
	return nil
}
