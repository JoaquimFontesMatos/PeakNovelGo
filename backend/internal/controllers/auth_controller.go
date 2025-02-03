package controllers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/validators"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	AuthService interfaces.AuthServiceInterface
	UserService interfaces.UserServiceInterface
}

func NewAuthController(authService interfaces.AuthServiceInterface, userService interfaces.UserServiceInterface) *AuthController {
	return &AuthController{
		AuthService: authService,
		UserService: userService,
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var registerRequest dtos.RegisterRequest

	if err := validators.ValidateBody(c, &registerRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the user service to create a new user
	if err := ac.UserService.RegisterUser(&registerRequest); err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.CONFLICT_ERROR:
				c.JSON(http.StatusConflict, gin.H{"error": myError.Message})
			}
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func (ac *AuthController) Login(c *gin.Context) {
	var req dtos.LoginRequest // Define a login request struct

	// Bind the incoming JSON to the LoginRequest model
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Validate credentials
	user, err := ac.AuthService.ValidateCredentials(req.Email, req.Password)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.INTERNAL_SERVER_ERROR:
				c.JSON(http.StatusInternalServerError, gin.H{"error": myError.Message})
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			default:
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			}
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT tokens (access token and refresh token)
	accessToken, refreshToken, err := ac.AuthService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIES_SECURE"))
	if err != nil {
		// Handle the error (e.g., log it or use a default value)
		secure = false // Default value if parsing fails
	}
	// Store refreshToken in HttpOnly cookie
	c.SetCookie(
		"refreshToken", // Name
		refreshToken,   // Value
		7*24*60*60,     // MaxAge: 7 days
		"/",            // Path
		"",             // Domain
		secure,         // Secure: Only send over HTTPS
		true,           // HttpOnly: Inaccessible to JavaScript
	)

	userDto, err := dtos.ConvertUserModelToDTO(*user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the access token and user info
	c.JSON(http.StatusOK, gin.H{
		"accessToken": accessToken,
		"user":        userDto,
	})
}

func (ac *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, user, err := ac.AuthService.RefreshToken(refreshToken)
	if err != nil {
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
				return
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
				return
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": userErr.Message})
				return
			}
		}

		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIES_SECURE"))
	if err != nil {
		// Handle the error (e.g., log it or use a default value)
		secure = false // Default value if parsing fails
	}

	// Send the new tokens back to the client
	ctx.SetCookie(
		"refreshToken",  // Name
		newRefreshToken, // Value
		7*24*60*60,      // MaxAge: 7 days
		"/",             // Path
		"",              // Domain
		secure,          // Secure: Only send over HTTPS
		true,            // HttpOnly: Inaccessible to JavaScript
	)

	userDto, err := dtos.ConvertUserModelToDTO(*user)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the access token and user info
	ctx.JSON(http.StatusOK, gin.H{
		"accessToken": newAccessToken,
		"user":        userDto,
	})
}

func (ac *AuthController) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	if err := ac.AuthService.Logout(refreshToken); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	secure, err := strconv.ParseBool(os.Getenv("COOKIES_SECURE"))
	if err != nil {
		// Handle the error (e.g., log it or use a default value)
		secure = false // Default value if parsing fails
	}

	ctx.SetCookie(
		"refreshToken",
		"",
		-1,
		"/",
		"",
		secure,
		true, // Set MaxAge to -1 to delete the cookie
	)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (ac *AuthController) VerifyEmail(c *gin.Context) {
	token := c.Query("token") // Extract the token from query parameters

	if err := ac.UserService.VerifyEmail(token); err != nil {
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case types.USER_NOT_FOUND_ERROR:
				c.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
				return
			case types.USER_DEACTIVATED_ERROR:
				c.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
				return
			case types.INVALID_TOKEN_ERROR:
				c.JSON(http.StatusUnauthorized, gin.H{"error": userErr.Message})
				return
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": userErr.Message})
				return
			}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
