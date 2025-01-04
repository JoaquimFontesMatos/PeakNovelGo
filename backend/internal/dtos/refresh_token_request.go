package dtos

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}