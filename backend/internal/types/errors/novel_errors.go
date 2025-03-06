package errors

import (
	"backend/internal/types"
	"net/http"
)

var (
	ErrAuthorRequired = &types.MyCustomError{
		Message:    "Author is required",
		StatusCode: http.StatusBadRequest,
		Code:       "AUTHOR_REQUIRED",
	}
	ErrAuthorTooLong = &types.MyCustomError{
		Message:    "Author cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       "AUTHOR_TOO_LONG",
	}
	ErrAuthorTooShort = &types.MyCustomError{
		Message:    "Author must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       "AUTHOR_TOO_SHORT",
	}
	ErrGenreRequired = &types.MyCustomError{
		Message:    "Genre is required",
		StatusCode: http.StatusBadRequest,
		Code:       "GENRE_REQUIRED",
	}
	ErrGenreTooLong = &types.MyCustomError{
		Message:    "Genre cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       "GENRE_TOO_LONG",
	}
	ErrGenreTooShort = &types.MyCustomError{
		Message:    "Genre must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       "GENRE_TOO_SHORT",
	}
	ErrTagRequired = &types.MyCustomError{
		Message:    "Tag is required",
		StatusCode: http.StatusBadRequest,
		Code:       "TAG_REQUIRED",
	}
	ErrTagTooLong = &types.MyCustomError{
		Message:    "Tag cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       "TAG_TOO_LONG",
	}
	ErrTagTooShort = &types.MyCustomError{
		Message:    "Tag must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       "TAG_TOO_SHORT",
	}
	ErrInvalidNovelUpdatesID = &types.MyCustomError{
		Message:    "Invalid novel updates ID specified (must be a string of lowercase letters, numbers, or single dashes, between 1 and 255 characters long)",
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_NOVEL_UPDATES_ID",
	}
	ErrInvalidLatestChapter = &types.MyCustomError{
		Message:    "Invalid latest chapter specified (must be an integer)",
		StatusCode: http.StatusBadRequest,
		Code:       "INVALID_LATEST_CHAPTER",
	}
	ErrNovelNotFound = &types.MyCustomError{
		Message:    "Novel not found",
		StatusCode: http.StatusNotFound,
		Code:       "NOVEL_NOT_FOUND",
	}
	ErrNovelAlreadyImported = &types.MyCustomError{
		Message:    "Novel already imported",
		StatusCode: http.StatusConflict,
		Code:       "NOVEL_ALREADY_IMPORTED",
	}
	ErrNovelConflict = &types.MyCustomError{
		Message:    "Novel conflict",
		StatusCode: http.StatusConflict,
		Code:       "NOVEL_CONFLICT",
	}
	ErrNoNovels = &types.MyCustomError{
		Message:    "No novels found",
		StatusCode: http.StatusNotFound,
		Code:       "NO_NOVELS",
	}
	ErrImportingNovel = &types.MyCustomError{
		Message:    "Failed to import novel",
		StatusCode: http.StatusInternalServerError,
		Code:       "IMPORTING_NOVEL",
	}
	ErrGettingNovels = &types.MyCustomError{
		Message:    "Failed to get novels",
		StatusCode: http.StatusInternalServerError,
		Code:       "GETTING_NOVELS",
	}
	ErrGettingNovel = &types.MyCustomError{
		Message:    "Failed to get novel",
		StatusCode: http.StatusInternalServerError,
		Code:       "GETTING_NOVELS",
	}
	ErrGettingTotalNovels = &types.MyCustomError{
		Message:    "Failed to get total number of novels",
		StatusCode: http.StatusInternalServerError,
		Code:       "GETTING_TOTAL_NOVELS",
	}
)

func HandleAssociationError(err error) error {
	return &types.MyCustomError{
		Message:    "Association error",
		StatusCode: http.StatusInternalServerError,
		Code:       "ASSOCIATION_ERROR",
		Wrapped:    err,
	}
}
