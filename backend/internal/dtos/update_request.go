package dtos

// UpdateRequest represents a request to update user information.
//
// Fields:
//   - Username (string): The user's username.
//   - Bio (string): The user's bio.
//   - ProfilePicture (string): The URL of the user's profile picture.
//   - PreferredLanguage (string): The user's preferred language.
//   - ReadingPreferences (string): The user's reading preferences.
//   - DateOfBirth (string): The user's date of birth.
//   - Roles (string): The user's roles.
type UpdateRequest struct {
	Username           string `json:"username,omitempty"`
	Bio                string `json:"bio,omitempty"`
	ProfilePicture     string `json:"profilePicture,omitempty"`
	PreferredLanguage  string `json:"preferredLanguage,omitempty"`
	ReadingPreferences string `json:"readingPreferences,omitempty"`
	DateOfBirth        string `json:"dateOfBirth,omitempty"`
	Roles              string `json:"roles,omitempty"`
}
