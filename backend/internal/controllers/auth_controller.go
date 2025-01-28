package controllers

import (
	"errors"
	"log"
	"net/http"

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

		log.Println(err)
		return
	}

	// Call the user service to create a new user
	if err := ac.UserService.RegisterUser(&registerRequest); err != nil {
		log.Println(err)

		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.CONFLICT_ERROR:
				c.JSON(http.StatusConflict, gin.H{"error": error.Message})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.INTERNAL_SERVER_ERROR:
				c.JSON(http.StatusInternalServerError, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				c.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
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

	// Store refreshToken in HttpOnly cookie
	c.SetCookie(
		"refreshToken", // Name
		refreshToken,   // Value
		7*24*60*60,     // MaxAge: 7 days
		"/",            // Path
		"",             // Domain
		false,          // Secure: Only send over HTTPS
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

func (c *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No refresh token provided"})
		return
	}

	// Call the service to refresh the token
	newAccessToken, newRefreshToken, user, err := c.AuthService.RefreshToken(refreshToken)
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

	// Send the new tokens back to the client
	ctx.SetCookie(
		"refreshToken",  // Name
		newRefreshToken, // Value
		7*24*60*60,      // MaxAge: 7 days
		"/",             // Path
		"",              // Domain
		false,           // Secure: Only send over HTTPS
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

	ctx.SetCookie(
		"refreshToken", "", -1, "/", "", true, true, // Set MaxAge to -1 to delete the cookie
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
