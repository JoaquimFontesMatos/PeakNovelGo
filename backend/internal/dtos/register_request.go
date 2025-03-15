package dtos

// RegisterRequest represents the request body for user registration.
//
// Fields:
//   - Username (string): The user's username. Must be provided.
//   - Email (string): The user's email address. Must be provided.
//   - Password (string): The user's password. Must be provided.
//   - Bio (string): The user's bio. Must be provided.
//   - ProfilePicture (string):  The URL of the user's profile picture (optional).
//   - DateOfBirth (string): The user's date of birth. Must be provided.
//   - Provider (string): The authentication provider (not included in JSON).
type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Bio            string `json:"bio" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
	DateOfBirth    string `json:"dateOfBirth" binding:"required"`
	Provider       string `json:"-"`
}
