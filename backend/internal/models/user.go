package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string    `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Email              string    `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password           string    `gorm:"size:255;not null" json:"-"`
	EmailVerified      bool      `gorm:"default:false" json:"emailVerified"`
	VerificationToken  string    `json:"-"`
	ProfilePicture     string    `gorm:"size:255" json:"profilePicture,omitempty"`
	Bio                string    `gorm:"size:500" json:"bio,omitempty"`
	Roles              string    `gorm:"size:255" json:"roles,omitempty"`
	LastLogin          time.Time `json:"lastLogin"`
	DateOfBirth        time.Time `json:"dateOfBirth"`
	PreferredLanguage  string    `gorm:"size:100" json:"preferredLanguage,omitempty"`
	ReadingPreferences string    `gorm:"size:255" json:"readingPreferences,omitempty"`
	IsDeleted          bool      `gorm:"default:false" json:"-"`
}
