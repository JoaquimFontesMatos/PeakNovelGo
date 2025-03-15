package models

import "time"

// RevokedToken represents a revoked access token.
//
// Fields:
//   - ID (uint): The unique identifier for the revoked token.
//   - Token (string): The revoked token string itself.  Must be unique.
//   - ExpiredAt (time.Time): The time when the token expired.
//   - CreatedAt (time.Time): The time the record was created (automatically updated by GORM).
//   - UpdatedAt (time.Time): The time the record was last updated (automatically updated by GORM).
type RevokedToken struct {
	ID uint `gorm:"primaryKey"`
	Token string `gorm:"uniqueIndex;not null"`
	ExpiredAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
