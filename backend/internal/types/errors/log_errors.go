package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	// Log errors
	LOGGING_ERROR = "LOGGING_ERROR"
	NO_LOGS_ERROR = "NO_LOGS_ERROR"

	// Log level errors
	LOG_LEVEL_NOT_FOUND_ERROR = "LOG_LEVEL_NOT_FOUND_ERROR"

	FETCHING_LOGS_ERROR = "FETCHING_LOGS_ERROR"
	FETCHING_LOG_ERROR  = "FETCHING_LOG_ERROR"
	TOTAL_LOGS_ERROR    = "TOTAL_LOGS_ERROR"
)

var (
	ErrNoLogs = &types.MyCustomError{
		Message:    "No logs found",
		StatusCode: http.StatusNotFound,
		Code:       NO_LOGS_ERROR,
	}
	ErrLogLevelNotFound = &types.MyCustomError{
		Message:    "Log level not found",
		StatusCode: http.StatusNotFound,
		Code:       LOG_LEVEL_NOT_FOUND_ERROR,
	}
)