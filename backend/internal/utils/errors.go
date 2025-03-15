package utils

import (
	"backend/internal/dtos"
	"backend/internal/types"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HandleError processes an error and sends an appropriate JSON response to the client based on the error type.
// If the error implements HTTPError, it uses its status and error code; otherwise, it returns a generic 500 error response.
//
// Parameters:
//   - c *gin.Context (context of the request)
//   - err error (error to be handled)
func HandleError(c *gin.Context, err error) {
	var httpErr types.HTTPError

	errorResponse := dtos.ErrorResponse{
		Error: "Internal server error",
		Code:  "INTERNAL_ERROR",
	}

	if errors.As(err, &httpErr) {
		errorResponse.Code = httpErr.ErrorCode()
		errorResponse.Error = httpErr.Error()
		c.JSON(httpErr.HTTPStatus(), errorResponse)
		return
	}

	c.JSON(http.StatusInternalServerError, errorResponse)
}
