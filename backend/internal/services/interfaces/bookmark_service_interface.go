package interfaces

import (
	"backend/internal/models"
)

type BookmarkServiceInterface interface {
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.Novel, int64, error)
	GetBookmarkByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error)
	UpdateBookmark(novel models.BookmarkedNovel) (models.BookmarkedNovel, error)
	CreateBookmark(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)
	UnbookmarkNovel(userID uint, novelID uint) error
}
