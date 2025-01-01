package validators

import (
	"github.com/gin-gonic/gin"
)

func ValidateParams(ctx *gin.Context) error {
	if len(ctx.Params) > 1 {
		return &ValidationError{Message: "Too many parameters"}
	} else if len(ctx.Params) == 0 {
		return &ValidationError{Message: "No parameters"}
	}

	return nil
}
