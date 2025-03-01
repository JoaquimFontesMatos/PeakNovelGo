package interfaces

import (
	"backend/internal/models"
)

type BookmarkRepositoryInterface interface {
	// GetBookmarkedNovelsByUserID gets the bookmarked novels by the given user
	//
	// Parameters:
	//   - userID uint (id of the user)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.Novel, int64, error)

	// GetBookmarkByUserIDAndNovelID gets a bookmarked novel by user ID and novel ID.
	//
	// Parameters:
	//   - userID uint (ID of the user)
	//   - novelID string (ID of the novel)
	//
	// Returns:
	//   - models.BookmarkedNovel (BookmarkedNovel struct)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be fetched
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be fetched
	GetBookmarkByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error)

	// UpdateBookmark updates a bookmarked novel in the database.
	//
	// Parameters:
	//   - novel models.BookmarkedNovel (BookmarkedNovel struct)
	//
	// Returns:
	//   - models.BookmarkedNovel (BookmarkedNovel struct)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be updated
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be updated
	UpdateBookmark(novel models.BookmarkedNovel) (models.BookmarkedNovel, error)

	// CreateBookmark creates a new bookmarked novel in the database.
	//
	// Parameters:
	//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
	//
	// Returns:
	//   - *models.BookmarkedNovel (pointer to BookmarkedNovel struct)
	//   - CONFLICT_ERROR if the bookmarked novel already exists
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be created
	CreateBookmark(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)

	// DeleteBookmark deletes a bookmarked novel from the database.
	//
	// Parameters:
	//   - userID uint (ID of the user)
	//   - novelID uint (ID of the novel)
	//
	// Returns:
	//   - error (nil if the bookmarked novel was deleted successfully, otherwise an error)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be deleted
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be deleted
	DeleteBookmark(userID uint, novelID uint) error
}
