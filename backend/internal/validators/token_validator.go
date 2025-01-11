package validators

import (
	"backend/internal/types"

	"time"
)

// IsVerificationTokenExpired checks if the verification token is expired.
//
// Parameters:
//   - createdAt time.Time (time when the token was created)
//   - emailVerified bool (whether the email has been verified)
//
// Returns:
//   - bool (true if the token is expired, false otherwise)
func IsVerificationTokenExpired(createdAt time.Time, emailVerified bool) bool {
	timeDifference := time.Since(createdAt)

	// Check if the token is expired
	if timeDifference >= 60*time.Minute || emailVerified {
		return true
	}
	return false
}

// ValidateToken validates the token and ensures it is valid.
//
// Parameters:
//   - token string (token)
//
// Returns:
//   - INVALID_TOKEN if the token is invalid
func ValidateToken(token string) error {
	if token == "" {
		return types.WrapError("INVALID_TOKEN_ERROR", "token is required", nil)
	}

	if len(token) > 255 {
		return types.WrapError("INVALID_TOKEN_ERROR", "token cannot be longer than 255 characters", nil)
	}

	return nil
}
