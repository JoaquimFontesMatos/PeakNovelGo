package validators

import (
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
