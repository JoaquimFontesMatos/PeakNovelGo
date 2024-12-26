package models

import "time"

// RevokedToken represents a revoked JWT token stored in the database.
type RevokedToken struct {
	ID         uint      `gorm:"primaryKey"`
	Token      string    `gorm:"uniqueIndex;not null"`
	ExpiredAt  time.Time `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}
