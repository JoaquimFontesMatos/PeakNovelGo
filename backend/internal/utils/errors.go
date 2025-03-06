package utils

import (
	"backend/internal/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// utils/errors.go
func HandleError(c *gin.Context, err error) {
	var httpErr types.HTTPError

	if errors.As(err, &httpErr) {
		c.JSON(httpErr.HTTPStatus(), gin.H{
			"error": httpErr.Error(),
			"code":  httpErr.ErrorCode(),
		})
		return
	}

	// Fallback for untyped errors
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
		"code":  "INTERNAL_ERROR",
	})
}
