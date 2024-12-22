package validators

import (
	"backend/internal/models"
	"errors"
	"regexp"
)

// ValidateUser checks the required fields in the User struct.
func ValidateUser(user *models.User) error {
	if user.Username == "" {
		return errors.New("username is required")
	}
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.Password == "" {
		return errors.New("password is required")
	}

	// Validate length constraints
	if len(user.Username) > 255 {
		return errors.New("username cannot be longer than 255 characters")
	}
	if len(user.Email) > 255 {
		return errors.New("email cannot be longer than 255 characters")
	}
	if len(user.Password) > 255 {
		return errors.New("password cannot be longer than 255 characters")
	}
	if len(user.ProfilePicture) > 255 {
		return errors.New("profile picture URL cannot be longer than 255 characters")
	}
	if len(user.Roles) > 255 {
		return errors.New("roles cannot be longer than 255 characters")
	}
	if len(user.PreferredLanguage) > 100 {
		return errors.New("preferred language cannot be longer than 100 characters")
	}
	if len(user.ReadingPreferences) > 255 {
		return errors.New("reading preferences cannot be longer than 255 characters")
	}
	if len(user.Bio) > 500 {
		return errors.New("bio cannot be longer than 500 characters")
	}

	// Optional: Validate email format using regex
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, user.Email)
	if err != nil {
		return errors.New("failed to validate email format")
	}
	if !matched {
		return errors.New("invalid email format")
	}

	return nil
}
