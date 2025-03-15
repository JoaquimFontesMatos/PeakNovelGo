package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	// Chapter errors
	CHAPTER_NOT_FOUND        = "CHAPTER_NOT_FOUND"
	CHAPTER_ALREADY_IMPORTED = "CHAPTER_ALREADY_IMPORTED"
	CHAPTER_CONFLICT         = "CHAPTER_CONFLICT"
	NO_CHAPTERS              = "NO_CHAPTERS"
	IMPORTING_CHAPTER        = "IMPORTING_CHAPTER"
	GETTING_CHAPTERS         = "GETTING_CHAPTERS"
	GETTING_CHAPTER          = "GETTING_CHAPTER"
	GETTING_TOTAL_CHAPTERS   = "GETTING_TOTAL_CHAPTERS"
)

var (
	ErrChapterNotFound = &types.MyCustomError{
		Message:    "Chapter not found",
		StatusCode: http.StatusNotFound,
		Code:       CHAPTER_NOT_FOUND,
	}
	ErrChapterAlreadyImported = &types.MyCustomError{
		Message:    "Chapter already imported",
		StatusCode: http.StatusConflict,
		Code:       CHAPTER_ALREADY_IMPORTED,
	}
	ErrChapterConflict = &types.MyCustomError{
		Message:    "Chapter conflict",
		StatusCode: http.StatusConflict,
		Code:       CHAPTER_CONFLICT,
	}
	ErrNoChapters = &types.MyCustomError{
		Message:    "No chapters found",
		StatusCode: http.StatusNotFound,
		Code:       NO_CHAPTERS,
	}
)
