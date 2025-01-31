package dtos

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Email          string `json:"email" binding:"required"`
	Password       string `json:"password" binding:"required"`
	Bio            string `json:"bio" binding:"required"`
	ProfilePicture string `json:"profilePicture"`
	DateOfBirth    string `json:"dateOfBirth" binding:"required"`
}
