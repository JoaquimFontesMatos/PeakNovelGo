package dtos

import (
	"backend/internal/models"
	"backend/internal/types"
	"encoding/json"
	"log"
	"time"

	"gorm.io/gorm"
)

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

// ConvertUserModelToDTO converts a User model to a UserDTO.
func ConvertUserModelToDTO(user models.User) (UserDTO, error) {
	var readingPreferences ReadingPreferences
	if user.ReadingPreferences != "" && user.ReadingPreferences != "null" {
		err := json.Unmarshal([]byte(user.ReadingPreferences), &readingPreferences)
		if err != nil {
			log.Printf("Error parsing ReadingPreferences: %s, Error: %v", user.ReadingPreferences, err)
			return UserDTO{}, types.WrapError(types.INTERNAL_SERVER_ERROR, "Error parsing reading preferences", err)
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
func ConvertUserDTOToModel(dto UserDTO) (models.User, error) {
	readingPreferencesJSON, err := json.Marshal(dto.ReadingPreferences)
	if err != nil {
		return models.User{}, types.WrapError(types.INTERNAL_SERVER_ERROR, "Error parsing reading preferences", err)
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
