package validators

import (
	"backend/internal/types"

	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ValidateBody validates the request body and ensures it is valid JSON.
//
// Parameters:
//   - ctx *gin.Context (gin context)
//   - body interface{} (body of the request)
//
// Returns:
//   - INVALID_BODY_ERROR if the request body is invalid JSON
func ValidateBody(ctx *gin.Context, body interface{}) error {
	// Read the raw body
	rawBody, err := ctx.GetRawData()
	if err != nil {
		return types.WrapError(types.INVALID_BODY_ERROR, "Failed to read request body", err)
	}

	// Check if it's valid JSON
	if !json.Valid(rawBody) {
		return types.WrapError(types.INVALID_BODY_ERROR, "Invalid JSON", nil)
	}

	// Reassign the raw body so ShouldBindJSON can read it
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Bind the JSON
	if err := ctx.ShouldBindJSON(body); err != nil {
		return types.WrapError(types.INVALID_BODY_ERROR, "Failed to bind JSON", err)
	}

	// Custom validation: ensure at least one field is non-empty
	if !hasAtLeastOneNonEmptyField(body) {
		return types.WrapError(types.INVALID_BODY_ERROR, "No fields provided", nil)
	}

	return nil
}

// hasAtLeastOneNonEmptyField checks if at least one field in a struct is non-empty
//
// Parameters:
//   - body interface{} (body of the request)
//
// Returns:
//   - bool (true if at least one field is non-empty, false otherwise)
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
