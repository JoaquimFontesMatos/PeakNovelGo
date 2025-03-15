package validators

import (
	"backend/internal/types"
	"backend/internal/types/errors"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"reflect"
)

// ValidateBody validates the request body. It checks if the body is valid JSON, binds it to the provided interface, and ensures at least one field is non-empty.
//
// Parameters:
//   - ctx (*gin.Context): The Gin context.
//   - body (interface{}): The interface to bind the request body to.
//
// Returns:
//   - error: An error (wrapped errors.INVALID_BODY_ERROR with status code http.StatusBadRequest) if the body is invalid or nil if the body is valid.
func ValidateBody(ctx *gin.Context, body interface{}) error {
	// Read the raw body
	rawBody, err := ctx.GetRawData()
	if err != nil {
		return types.WrapError(errors.INVALID_BODY_ERROR, "Failed to read request body", http.StatusBadRequest, err)
	}

	// Check if it's valid JSON
	if !json.Valid(rawBody) {
		return types.WrapError(errors.INVALID_BODY_ERROR, "Invalid JSON", http.StatusBadRequest, nil)
	}

	// Reassign the raw body so ShouldBindJSON can read it
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Bind the JSON
	if err := ctx.ShouldBindJSON(body); err != nil {
		return types.WrapError(errors.INVALID_BODY_ERROR, "Failed to bind JSON", http.StatusBadRequest, err)
	}

	// Custom validation: ensure at least one field is non-empty
	if !hasAtLeastOneNonEmptyField(body) {
		return types.WrapError(errors.INVALID_BODY_ERROR, "No fields provided", http.StatusBadRequest, nil)
	}

	return nil
}

// hasAtLeastOneNonEmptyField checks if at least one field in a struct is non-empty
//
// Parameters:
//   - body (interface{}): body of the request
//
// Returns:
//   - bool: true if at least one field is non-empty, false otherwise
func hasAtLeastOneNonEmptyField(body interface{}) bool {
	v := reflect.ValueOf(body).Elem() // Dereference the pointer to access the struct
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field is non-empty
		if !field.IsZero() {
			return true
		}
	}
	return false
}
