package dtos

type UpdateRequest struct {
    Username           string `json:"username,omitempty"`
    Bio                string `json:"bio,omitempty"`
    ProfilePicture     string `json:"profilePicture,omitempty"`
    PreferredLanguage  string `json:"preferredLanguage,omitempty"`
    ReadingPreferences string `json:"readingPreferences,omitempty"`
    DateOfBirth        string `json:"dateOfBirth,omitempty"`
    Roles              string `json:"roles,omitempty"`
}
