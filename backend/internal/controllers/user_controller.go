package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/validators"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service interfaces.UserServiceInterface
}

func NewUserController(service interfaces.UserServiceInterface) *UserController {
	return &UserController{service: service}
}

// HandleGetUser handles GET /users/:id
func (c *UserController) HandleGetUser(ctx *gin.Context) {
	if len(ctx.Params) > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Too many parameters"})
		return
	} else if len(ctx.Params) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No parameters"})
		return
	}

	// Get the user ID from the URL parameters
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
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
	if len(ctx.Params) > 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Too many parameters"})
		return
	} else if len(ctx.Params) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No parameters"})
		return
	}

	email := ctx.Param("email")

	// Call the service layer to fetch the user
	user, err := c.service.GetUserByEmail(email)
	if err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleGetUserByUsername handles GET /users/username/:username
func (c *UserController) HandleGetUserByUsername(ctx *gin.Context) {
	if err := validators.ValidateParams(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	username := ctx.Param("username")

	// Call the service layer to fetch the user
	user, err := c.service.GetUserByUsername(username)
	if err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleUpdateUser handles PUT /users/:i
func (c *UserController) UpdateUserFields(ctx *gin.Context) {
	var updateFields models.UpdateFields

	if err := validators.ValidateBody(ctx, &updateFields); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	if err := c.service.UpdateUserFields(uid, updateFields); err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var userErr *validators.UserError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case "USER_NOT_FOUND":
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
			case "USER_DELETED":
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func (c *UserController) UpdatePassword(ctx *gin.Context) {
	var req struct {
		CurrentPassword string `json:"current_password" binding:"required"`
		NewPassword     string `json:"new_password" binding:"required"`
	}

	if err := validators.ValidateBody(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	if err := c.service.UpdatePassword(uid, req.CurrentPassword, req.NewPassword); err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var userErr *validators.UserError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case "USER_NOT_FOUND":
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
			case "USER_DELETED":
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
			case "PASSWORD_DIFF":
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": userErr.Message})
			case "INVALID_PASSWORD":
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": userErr.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}

func (c *UserController) UpdateEmail(ctx *gin.Context) {
	var req struct {
		NewEmail string `json:"new_email" binding:"required,email"`
	}

	// Validate request body
	if err := validators.ValidateBody(ctx, &req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	if err := c.service.UpdateEmail(uid, req.NewEmail); err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var userErr *validators.UserError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case "USER_NOT_FOUND":
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
			case "USER_DELETED":
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Email update initiated. Please verify the new email."})
}

// HandleDeleteUser handles DELETE /users/:id
func (c *UserController) HandleDeleteUser(ctx *gin.Context) {
	if err := validators.ValidateParams(ctx); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	uid := uint(id)

	// Call the service layer to delete the user
	err = c.service.DeleteUser(uid)
	if err != nil {
		if _, ok := err.(*validators.ValidationError); ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
