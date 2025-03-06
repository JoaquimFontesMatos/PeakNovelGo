package errors

import (
	"backend/internal/types"
	"net/http"
)

var (
	ErrScriptNotFound = &types.MyCustomError{
		Message:    "Script not found",
		StatusCode: http.StatusNotFound,
		Code:       "SCRIPT_NOT_FOUND",
	}
	ErrScriptAlreadyImported = &types.MyCustomError{
		Message:    "Script already imported",
		StatusCode: http.StatusConflict,
		Code:       "SCRIPT_ALREADY_IMPORTED",
	}
	ErrScriptConflict = &types.MyCustomError{
		Message:    "Script conflict",
		StatusCode: http.StatusConflict,
		Code:       "SCRIPT_CONFLICT",
	}
	ErrNoScript = &types.MyCustomError{
		Message:    "No script found",
		StatusCode: http.StatusNotFound,
		Code:       "NO_SCRIPT",
	}
)
