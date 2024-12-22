package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username           string    `gorm:"size:255;unique;not null" json:"username"`
	Email              string    `gorm:"size:255;unique;not null" json:"email"`
	Password           string    `gorm:"size:255;not null" json:"password"`
	EmailVerified      bool      `json:"email_verified"`
	VerificationToken  string    `json:"verification_token"`
	ProfilePicture     string    `gorm:"size:255" json:"profile_picture"`
	Bio                string    `gorm:"size:500" json:"bio"`
	Roles              string    `gorm:"size:255" json:"roles"`
	LastLogin          time.Time `json:"last_login"`
	DateOfBirth        time.Time `json:"date_of_birth"`
	PreferredLanguage  string    `gorm:"size:100" json:"preferred_language"`
	ReadingPreferences string    `gorm:"size:255" json:"reading_preferences"`
	IsDeleted          bool      `gorm:"default:false" json:"is_deleted"`
}
