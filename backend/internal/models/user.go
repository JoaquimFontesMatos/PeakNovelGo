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
	EmailVerified      bool      `gorm:"default:false" json:"email_verified"`
	VerificationToken  string    `json:"-"`
	ProfilePicture     string    `gorm:"size:255" json:"profile_picture,omitempty"`
	Bio                string    `gorm:"size:500" json:"bio,omitempty"`
	Roles              string    `gorm:"size:255" json:"roles,omitempty"`
	LastLogin          time.Time `json:"last_login"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	PreferredLanguage  string    `gorm:"size:100" json:"preferred_language,omitempty"`
	ReadingPreferences string    `gorm:"size:255" json:"reading_preferences,omitempty"`
	IsDeleted          bool      `gorm:"default:false" json:"is_deleted"`
}
