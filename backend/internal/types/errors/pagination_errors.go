package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	INVALID_PAGE      = "INVALID_PAGE"
	INVALID_LIMIT     = "INVALID_LIMIT"
	PAGE_OUT_OF_RANGE = "PAGE_OUT_OF_RANGE"
	NO_RESULTS        = "NO_RESULTS"
)

var (
	ErrInvalidPage = &types.MyCustomError{
		Message:    "Invalid page specified (must be an integer, between 1 and 1000)",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_PAGE,
	}
	ErrInvalidLimit = &types.MyCustomError{
		Message:    "Invalid limit specified (must be an integer, between 10 and 100)",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_LIMIT,
	}
	ErrPageOutOfRange = &types.MyCustomError{
		Message:    "Page out of range",
		StatusCode: http.StatusBadRequest,
		Code:       PAGE_OUT_OF_RANGE,
	}
	ErrNoResults = &types.MyCustomError{
		Message:    "No results found",
		StatusCode: http.StatusNotFound,
		Code:       NO_RESULTS,
	}
)
