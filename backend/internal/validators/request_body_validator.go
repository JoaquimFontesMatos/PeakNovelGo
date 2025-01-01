package validators

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"

	"github.com/gin-gonic/gin"
)

func ValidateBody(ctx *gin.Context, body interface{}) error {
	// Read the raw body
	rawBody, err := ctx.GetRawData()
	if err != nil {
		return &ValidationError{Message: "Invalid Input"}
	}

	// Check if it's valid JSON
	if !json.Valid(rawBody) {
		return &ValidationError{Message: "Invalid Input"}
	}

	// Reassign the raw body so ShouldBindJSON can read it
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	// Bind the JSON
	if err := ctx.ShouldBindJSON(body); err != nil {
		return &ValidationError{Message: "Invalid Input"}
	}

	// Custom validation: ensure at least one field is non-empty
	if !hasAtLeastOneNonEmptyField(body) {
		return &ValidationError{Message: "At least one field must be provided"}
	}

	return nil
}

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