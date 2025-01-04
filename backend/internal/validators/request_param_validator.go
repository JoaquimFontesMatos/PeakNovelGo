package validators

import (
	"backend/internal/types"
	"github.com/gin-gonic/gin"
)

func ValidateParams(ctx *gin.Context) error {
	if len(ctx.Params) > 1 {
		return &types.ValidationError{Message: "Too many parameters"}
	} else if len(ctx.Params) == 0 {
		return &types.ValidationError{Message: "No parameters"}
	}

	return nil
}
