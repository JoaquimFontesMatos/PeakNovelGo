package validators

import (
	"backend/internal/models"
	"errors"
	"regexp"
	"time"
)

func ValidateUser(user *models.User) error {
	if err := ValidateEmail(user.Email); err != nil {
		return err
	}

	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	fields := models.UpdateFields{
		Bio:                user.Bio,
		PreferredLanguage:  user.PreferredLanguage,
		Roles:              user.Roles,
		Username:           user.Username,
		ProfilePicture:     user.ProfilePicture,
		ReadingPreferences: user.ReadingPreferences,
		DateOfBirth:        user.DateOfBirth.Format("2006-01-02"),
	}

	if err := ValidateUserFields(fields); err != nil {
		return err
	}

	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 255 {
		return errors.New("password cannot be longer than 255 characters")
	}
	return nil
}

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}

	if len(email) > 255 {
		return errors.New("email cannot be longer than 255 characters")
	}

	// Optional: Validate email format using regex
	emailRegex := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, err := regexp.MatchString(emailRegex, email)
	if err != nil {
		return errors.New("failed to validate email format")
	}
	if !matched {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) > 255 {
		return errors.New("username cannot be longer than 255 characters")
	}
	return nil
}

// ValidateUserFields checks the fields being updated and ensures they meet validation criteria.
// ValidateUserFields validates an UpdateFields struct.
func ValidateUserFields(fields models.UpdateFields) error {
	// Check the username length
	if len(fields.Username) > 255 {
		return errors.New("username cannot be longer than 255 characters")
	}

	// Check the bio length
	if len(fields.Bio) > 500 {
		return errors.New("bio cannot be longer than 500 characters")
	}

	// Check the profile picture URL length
	if len(fields.ProfilePicture) > 255 {
		return errors.New("profile picture URL cannot be longer than 255 characters")
	}

	// Check the preferred language length
	if len(fields.PreferredLanguage) > 100 {
		return errors.New("preferred language cannot be longer than 100 characters")
	}

	// Check the reading preferences length
	if len(fields.ReadingPreferences) > 255 {
		return errors.New("reading preferences cannot be longer than 255 characters")
	}

	// Check the roles length
	if len(fields.Roles) > 255 {
		return errors.New("roles cannot be longer than 255 characters")
	}

	// Validate the date of birth
	if fields.DateOfBirth != "" { // Skip validation if the field is empty
		dob, err := time.Parse("2006-01-02", fields.DateOfBirth)
		if err != nil {
			return errors.New("date of birth must be a valid date in YYYY-MM-DD format")
		}
		if dob.After(time.Now().AddDate(-18, 0, 0)) {
			return errors.New("you must be at least 18 years old")
		}
	}

	return nil
}
