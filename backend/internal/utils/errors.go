package utils

import (
	"backend/internal/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError processes an error and sends an appropriate JSON response to the client based on the error type.
// If the error implements HTTPError, it uses its status and error code; otherwise, it returns a generic 500 error response.
//
// Parameters:
// 	- c *gin.Context (context of the request)
// 	- err error (error to be handled)
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
