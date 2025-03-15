package dtos

// LoginRequest represents the request body for user login.
//
// Fields:
//   - Email (string): The user's email address. Must be provided.
//   - Password (string): The user's password. Must be provided.
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
