package models

type UpdateFields struct {
    Username          string `json:"username,omitempty"`
    Bio               string `json:"bio,omitempty"`
    ProfilePicture    string `json:"profile_picture,omitempty"`
    PreferredLanguage string `json:"preferred_language,omitempty"`
    ReadingPreferences string `json:"reading_preferences,omitempty"`
    DateOfBirth       string `json:"date_of_birth,omitempty"`
    Roles             string `json:"roles,omitempty"`
}
