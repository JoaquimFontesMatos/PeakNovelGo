package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	// Author errors
	INVALID_AUTHOR=  "INVALID_AUTHOR"
	AUTHOR_ASSOCIATION_ERROR = "AUTHOR_ASSOCIATION_ERROR"

	// Genre errors
	INVALID_GENRE = "INVALID_GENRE"
	GENRE_ASSOCIATION_ERROR = "GENRE_ASSOCIATION_ERROR"

	// Tag errors
	INVALID_TAG = "INVALID_TAG"
	TAG_ASSOCIATION_ERROR = "TAG_ASSOCIATION_ERROR"

	// NovelUpdatesId errors
	INVALID_NOVEL_UPDATES_ID

	// LatestChapter errors
	INVALID_LATEST_CHAPTER

	// Novel errors
	NOVEL_NOT_FOUND
	NOVEL_ALREADY_IMPORTED
	NOVEL_CONFLICT
	NO_NOVELS
	IMPORTING_NOVEL
	GETTING_NOVELS
	GETTING_NOVEL
	GETTING_TOTAL_NOVELS
)

var (
	ErrAuthorRequired = &types.MyCustomError{
		Message:    "Author is required",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_AUTHOR,
	}
	ErrAuthorTooLong = &types.MyCustomError{
		Message:    "Author cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_AUTHOR,
	}
	ErrAuthorTooShort = &types.MyCustomError{
		Message:    "Author must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_AUTHOR,
	}
	ErrGenreRequired = &types.MyCustomError{
		Message:    "Genre is required",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_GENRE,
	}
	ErrGenreTooLong = &types.MyCustomError{
		Message:    "Genre cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_GENRE,
	}
	ErrGenreTooShort = &types.MyCustomError{
		Message:    "Genre must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_GENRE,
	}
	ErrTagRequired = &types.MyCustomError{
		Message:    "Tag is required",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_TAG,
	}
	ErrTagTooLong = &types.MyCustomError{
		Message:    "Tag cannot be longer than 255 characters",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_TAG,
	}
	ErrTagTooShort = &types.MyCustomError{
		Message:    "Tag must be at least 1 characters long",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_TAG,
	}
	ErrInvalidNovelUpdatesID = &types.MyCustomError{
		Message:    "Invalid novel updates ID specified (must be a string of lowercase letters, numbers, or single dashes, between 1 and 255 characters long)",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_NOVEL_UPDATES_ID,
	}
	ErrInvalidLatestChapter = &types.MyCustomError{
		Message:    "Invalid latest chapter specified (must be an integer)",
		StatusCode: http.StatusBadRequest,
		Code:       INVALID_LATEST_CHAPTER,
	}
	ErrNovelNotFound = &types.MyCustomError{
		Message:    "Novel not found",
		StatusCode: http.StatusNotFound,
		Code:       NOVEL_NOT_FOUND,
	}
	ErrNovelAlreadyImported = &types.MyCustomError{
		Message:    "Novel already imported",
		StatusCode: http.StatusConflict,
		Code:       NOVEL_ALREADY_IMPORTED,
	}
	ErrNovelConflict = &types.MyCustomError{
		Message:    "Novel conflict",
		StatusCode: http.StatusConflict,
		Code:       NOVEL_CONFLICT,
	}
	ErrNoNovels = &types.MyCustomError{
		Message:    "No novels found",
		StatusCode: http.StatusNotFound,
		Code:       NO_NOVELS,
	}
	ErrImportingNovel = &types.MyCustomError{
		Message:    "Failed to import novel",
		StatusCode: http.StatusInternalServerError,
		Code:       IMPORTING_NOVEL,
	}
	ErrGettingNovels = &types.MyCustomError{
		Message:    "Failed to get novels",
		StatusCode: http.StatusInternalServerError,
		Code:       GETTING_NOVELS,
	}
	ErrGettingNovel = &types.MyCustomError{
		Message:    "Failed to get novel",
		StatusCode: http.StatusInternalServerError,
		Code:       GETTING_NOVELS,
	}
	ErrGettingTotalNovels = &types.MyCustomError{
		Message:    "Failed to get total number of novels",
		StatusCode: http.StatusInternalServerError,
		Code:       GETTING_TOTAL_NOVELS,
	}
)
