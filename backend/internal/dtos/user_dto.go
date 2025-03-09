package dtos

import (
	"backend/internal/models"
	"backend/internal/types"
	"backend/internal/types/errors"
	"encoding/json"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// UserDTO represents a user data transfer object.
//
// Fields:
//   - Username (string): The username of the user.
//   - Email (string): The email address of the user.
//   - EmailVerified (bool): Indicates whether the user's email address is verified.
//   - ProfilePicture (string): The URL of the user's profile picture (optional).
//   - Bio (string): A short biography of the user (optional).
//   - Roles (string): A string representing the user's roles (optional).
//   - LastLogin (time.Time): The timestamp of the user's last login.
//   - DateOfBirth (time.Time): The user's date of birth.
//   - PreferredLanguage (string): The user's preferred language (optional).
//   - ReadingPreferences (ReadingPreferences): The user's reading preferences (optional).
type UserDTO struct {
	gorm.Model
	Username           string             `json:"username"`
	Email              string             `json:"email"`
	EmailVerified      bool               `json:"emailVerified"`
	ProfilePicture     string             `json:"profilePicture,omitempty"`
	Bio                string             `json:"bio,omitempty"`
	Roles              string             `json:"roles,omitempty"`
	LastLogin          time.Time          `json:"lastLogin"`
	DateOfBirth        time.Time          `json:"dateOfBirth"`
	PreferredLanguage  string             `json:"preferredLanguage,omitempty"`
	ReadingPreferences ReadingPreferences `json:"readingPreferences,omitempty"`
}

// ConvertUserModelToDTO converts a models.User to a UserDTO.  Handles potential errors during JSON unmarshalling of
// ReadingPreferences.
//
// Parameters:
//   - user (models.User): The user model to convert.
//
// Returns:
//   - UserDTO: The converted UserDTO.
//   - error: An error occurred during conversion, specifically if the ReadingPreferences field cannot be unmarshalled.
//     The error type will be either types.Error or an underlying error from json.Unmarshal.
//
// Error types:
//   - types.Error: Wrapped error indicating invalid reading preferences. The original error will be included in the wrapped
//     error.
func ConvertUserModelToDTO(user models.User) (UserDTO, error) {
	var readingPreferences ReadingPreferences
	if user.ReadingPreferences != "" && user.ReadingPreferences != "null" {
		err := json.Unmarshal([]byte(user.ReadingPreferences), &readingPreferences)
		if err != nil {
			log.Printf("Error parsing ReadingPreferences: %s, Error: %v", user.ReadingPreferences, err)
			return UserDTO{}, types.WrapError(errors.INVALID_READING_PREFERENCES, "Error parsing reading preferences", http.StatusBadRequest, err)
		}
	}

	return UserDTO{
		Model: gorm.Model{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			DeletedAt: user.DeletedAt,
		},
		Username:           user.Username,
		Email:              user.Email,
		EmailVerified:      user.EmailVerified,
		ProfilePicture:     user.ProfilePicture,
		Bio:                user.Bio,
		Roles:              user.Roles,
		LastLogin:          user.LastLogin,
		DateOfBirth:        user.DateOfBirth,
		PreferredLanguage:  user.PreferredLanguage,
		ReadingPreferences: readingPreferences,
	}, nil
}

// ConvertUserDTOToModel converts a UserDTO to a User model.
//
// Parameters:
//   - dto (UserDTO): The UserDTO to convert.
//
// Returns:
//   - models.User: The converted User model.
//   - error: An error if the conversion fails.  This may include errors related to parsing reading preferences.
//
// Error types:
//   - types.Error:  Wraps underlying errors.  Specifically, errors with code `errors.INVALID_READING_PREFERENCES` will
//     be returned if the reading preferences JSON cannot be marshaled.  The HTTP status code will be 400 (Bad Request).
func ConvertUserDTOToModel(dto UserDTO) (models.User, error) {
	readingPreferencesJSON, err := json.Marshal(dto.ReadingPreferences)
	if err != nil {
		return models.User{}, types.WrapError(errors.INVALID_READING_PREFERENCES, "Error parsing reading preferences", http.StatusBadRequest, err)
	}

	return models.User{
		Model: gorm.Model{
			ID:        dto.ID,
			CreatedAt: dto.CreatedAt,
			UpdatedAt: dto.UpdatedAt,
			DeletedAt: dto.DeletedAt,
		},
		Username:           dto.Username,
		Email:              dto.Email,
		EmailVerified:      dto.EmailVerified,
		ProfilePicture:     dto.ProfilePicture,
		Bio:                dto.Bio,
		Roles:              dto.Roles,
		LastLogin:          dto.LastLogin,
		DateOfBirth:        dto.DateOfBirth,
		PreferredLanguage:  dto.PreferredLanguage,
		ReadingPreferences: string(readingPreferencesJSON),
	}, nil
}
