package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/types"
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the user
	ctx.JSON(http.StatusOK, user)
}

// HandleUpdateUser handles PUT /users/:i
func (c *UserController) UpdateUserFields(ctx *gin.Context) {
	var updateFields dtos.UpdateRequest

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
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			switch userErr.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": userErr.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": userErr.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": userErr.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			case types.PASSWORD_DIFF_ERROR:
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": error.Message})
			case types.INVALID_PASSWORD_ERROR:
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.INVALID_CREDENTIALS_ERROR:
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": error.Message})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
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
		var error *types.MyError
		if errors.As(err, &error) {
			switch error.Code {
			case types.USER_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": error.Message})
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": error.Message})
			case types.USER_DEACTIVATED_ERROR:
				ctx.JSON(http.StatusForbidden, gin.H{"error": error.Message})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
