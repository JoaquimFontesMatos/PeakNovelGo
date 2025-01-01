package models

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Bio      string `json:"bio" binding:"required"`
	ProfilePicture     string `json:"profile_picture" binding:"required"`
	DateOfBirth        string `json:"date_of_birth" binding:"required"`
}