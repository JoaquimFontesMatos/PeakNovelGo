package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system.  It includes user details, login information, and preferences.
//
// Parameters:
//   - Username (string): The user's username. Must be unique and not null. Maximum length is 255 characters.
//   - Email (string): The user's email address. Must be unique and not null. Maximum length is 255 characters.
//   - Password (string): The user's password. Not null. This field is not exposed in JSON responses.
//   - EmailVerified (bool): Indicates if the user's email address has been verified. Defaults to false.
//   - VerificationToken (string): A token used for email verification. This field is not exposed in JSON responses.
//   - ProfilePicture (string): URL to the user's profile picture. Maximum length is 255 characters. Optional.
//   - Bio (string): A short biography of the user. Maximum length is 500 characters. Optional.
//   - Roles (string): A comma-separated string of roles assigned to the user. Maximum length is 255 characters. Optional.
//   - LastLogin (time.Time): The timestamp of the user's last login.
//   - DateOfBirth (time.Time): The user's date of birth.
//   - PreferredLanguage (string): The user's preferred language. Maximum length is 100 characters. Optional.
//   - ReadingPreferences (string): The user's reading preferences. Maximum length is 255 characters. Optional.
//   - IsDeleted (bool): Indicates if the user's account has been deleted. Defaults to false. This field is not exposed
//     in JSON responses.
//   - Provider (string): The authentication provider used to create the user account. Maximum length is 255 characters.
//     This field is not exposed in JSON responses.
type User struct {
	gorm.Model
	Username string `gorm:"size:255;uniqueIndex;not null" json:"username"`
	Email string `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password string `gorm:"size:255;not null" json:"-"`
	EmailVerified bool `gorm:"default:false" json:"emailVerified"`
	VerificationToken string `json:"-"`
	ProfilePicture string `gorm:"size:255" json:"profilePicture,omitempty"`
	Bio string `gorm:"size:500" json:"bio,omitempty"`
	Roles string `gorm:"size:255" json:"roles,omitempty"`
	LastLogin time.Time `json:"lastLogin"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	PreferredLanguage string `gorm:"size:100" json:"preferredLanguage,omitempty"`
	ReadingPreferences string `gorm:"size:255" json:"readingPreferences,omitempty"`
	IsDeleted bool `gorm:"default:false" json:"-"`
	Provider string `gorm:"size:255" json:"-"`
}
