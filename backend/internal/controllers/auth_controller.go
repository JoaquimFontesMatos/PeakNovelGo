package controllers

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/services"
	"github.com/gin-gonic/gin" // Assuming you're using Gin framework
)

type AuthController struct {
	AuthService services.AuthService
	UserService services.UserService
}

func NewAuthController(authService services.AuthService, userService services.UserService) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req models.User // Define a request struct in `models`

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the user service to create a new user
	if err := ac.UserService.RegisterUser(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req models.LoginRequest // Define a login request struct

	// Bind the incoming JSON to the LoginRequest model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate credentials
	user, err := ac.AuthService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT tokens (access token and refresh token)
	accessToken, refreshToken, err := ac.AuthService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	// Respond with both tokens
	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	// Extract the refresh token from the request header
	refreshToken := ctx.GetHeader("Authorization")
	if refreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, err := c.AuthService.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Send the new tokens back to the client
	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	refreshToken := ctx.GetHeader("Authorization")
	if refreshToken == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token is required"})
		return
	}

	if err := ac.AuthService.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token") // Extract the token from query parameters

	if err := ac.UserService.VerifyEmail(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
