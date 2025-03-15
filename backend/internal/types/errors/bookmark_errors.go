package errors

import (
	"backend/internal/types"
	"net/http"
)

const (
	// Bookmark errors
	BOOKMARK_NOT_FOUND        = "BOOKMARK_NOT_FOUND"
	BOOKMARK_ALREADY_IMPORTED = "BOOKMARK_ALREADY_IMPORTED"
	BOOKMARK_CONFLICT         = "BOOKMARK_CONFLICT"
	NO_BOOKMARKS              = "NO_BOOKMARKS"
	IMPORTING_BOOKMARK        = "IMPORTING_BOOKMARK"
	GETTING_BOOKMARKS         = "GETTING_BOOKMARKS"
	GETTING_BOOKMARK          = "GETTING_BOOKMARK"
	GETTING_TOTAL_BOOKMARKS   = "GETTING_TOTAL_BOOKMARKS"
	UPDATING_BOOKMARK         = "UPDATING_BOOKMARK"
	DELETING_BOOKMARK         = "DELETING_BOOKMARK"
)

var (
	ErrBookmarkNotFound = &types.MyCustomError{
		Message:    "Bookmark not found",
		StatusCode: http.StatusNotFound,
		Code:       BOOKMARK_NOT_FOUND,
	}
	ErrBookmarkAlreadyImported = &types.MyCustomError{
		Message:    "Bookmark already imported",
		StatusCode: http.StatusConflict,
		Code:       BOOKMARK_ALREADY_IMPORTED,
	}
	ErrBookmarkConflict = &types.MyCustomError{
		Message:    "Bookmark conflict",
		StatusCode: http.StatusConflict,
		Code:       BOOKMARK_CONFLICT,
	}
	ErrNoBookmarks = &types.MyCustomError{
		Message:    "No bookmarks found",
		StatusCode: http.StatusNotFound,
		Code:       NO_BOOKMARKS,
	}
)
