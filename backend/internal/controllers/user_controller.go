package controllers

import (
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/services/interfaces"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service interfaces.UserServiceInterface
}

func NewUserController(service interfaces.UserServiceInterface) *UserController {
	return &UserController{service: service}
}

// HandleCreateUser handles POST /users
func (c *UserController) HandleCreateUser(ctx *gin.Context) {
	var user models.User
	// Bind JSON input to the user struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to register the user
	if err := c.service.RegisterUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	// Respond with the created user
	ctx.JSON(http.StatusCreated, user)
}

// HandleGetUser handles GET /users/:id
func (c *UserController) HandleGetUser(ctx *gin.Context) {
	// Get the user ID from the URL parameters
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	// Call the service layer to fetch the user
	user, err := c.service.GetUser(uid)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleGetUserByEmail handles GET /users/email/:email
func (c *UserController) HandleGetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	// Call the service layer to fetch the user
	user, err := c.service.GetUserByEmail(email)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleGetUserByUsername handles GET /users/username/:username
func (c *UserController) HandleGetUserByUsername(ctx *gin.Context) {
	username := ctx.Param("username")

	// Call the service layer to fetch the user
	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleUpdateUser handles PUT /users/:id
func (c *UserController) HandleUpdateUser(ctx *gin.Context) {
	var user models.User
	// Bind JSON input to the user struct
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to update the user
	if err := c.service.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Respond with the updated user
	ctx.JSON(http.StatusOK, user)
}

// HandleDeleteUser handles DELETE /users/:id
func (c *UserController) HandleDeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	uid := uint(id)

	// Call the service layer to delete the user
	err = c.service.DeleteUser(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// HandleVerifyEmail handles POST /users/verify-email
func (c *UserController) HandleVerifyEmail(ctx *gin.Context) {
	var verificationToken string
	if err := ctx.ShouldBindJSON(&verificationToken); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Call the service layer to verify the email
	err := c.service.VerifyEmail(verificationToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Email verified successfully"})
}
