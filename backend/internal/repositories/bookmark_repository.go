package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"log"

	"gorm.io/gorm"
)

// BookmarkRepository struct represents a repository for chapter management.
type BookmarkRepository struct {
	db *gorm.DB
}

// NewBookmarkRepository creates a new BookmarkRepository instance
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *BookmarkRepository (pointer to the BookmarkRepository instance)
func NewBookmarkRepository(db *gorm.DB) *BookmarkRepository {
	return &BookmarkRepository{db: db}
}

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
func (b *BookmarkRepository) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the user
	if err := b.db.Model(&models.Novel{}).
		Joins("JOIN bookmarked_novels ON bookmarked_novels.novel_id = novels.id").
		Where("bookmarked_novels.user_id = ?", userID).
		Count(&total).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of bookmarked novels", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := b.db.Model(&models.Novel{}).
		Joins("JOIN bookmarked_novels ON bookmarked_novels.novel_id = novels.id").
		Where("bookmarked_novels.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novels", err)
	}

	return novels, total, nil
}

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
func (b *BookmarkRepository) GetBookmarkByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error) {
	var novel models.BookmarkedNovel
	if err := b.db.Model(&models.BookmarkedNovel{}).
		Joins("JOIN novels ON novels.id = bookmarked_novels.novel_id").
		Where("bookmarked_novels.user_id = ? AND novels.novel_updates_id = ?", userID, novelID).
		First(&novel).Error; err != nil {

		if err.Error() == "record not found" {
			return novel, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return novel, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch bookmarked novel", err)
	}
	return novel, nil
}

// UpdateBookmark updates a bookmarked novel in the database.
//
// Parameters:
//   - novel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - models.BookmarkedNovel (BookmarkedNovel struct)
//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be updated
//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be updated
func (b *BookmarkRepository) UpdateBookmark(novel models.BookmarkedNovel) (models.BookmarkedNovel, error) {
	if err := b.db.Model(&models.BookmarkedNovel{}).
		Where("user_id = ? AND novel_id = ?", novel.UserID, novel.NovelID).
		Update("status", novel.Status).
		Update("score", novel.Score).
		Update("current_chapter", novel.CurrentChapter).
		Error; err != nil {

		if err.Error() == "record not found" {
			return novel, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return novel, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to update bookmarked novel", err)
	}
	return novel, nil
}

// CreateBookmark creates a new bookmarked novel in the database.
//
// Parameters:
//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - *models.BookmarkedNovel (pointer to BookmarkedNovel struct)
//   - CONFLICT_ERROR if the bookmarked novel already exists
//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be created
func (b *BookmarkRepository) CreateBookmark(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error) {
	if IsBookmarkedNovelCreated := b.isBookmarkCreated(bookmarkedNovel); IsBookmarkedNovelCreated {
		return nil, types.WrapError(types.CONFLICT_ERROR, "Bookmarked novel already exists", nil)
	}

	// Save the bookmarked novel
	if err := b.db.Create(&bookmarkedNovel).Error; err != nil {
		log.Println(err)
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to create bookmarked novel", err)
	}

	return &bookmarkedNovel, nil
}

// isBookmarkCreated checks if a bookmarked novel with the given novel ID and user ID already exists in the database.
//
// Parameters:
//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - bool (true if the bookmarked novel already exists, false otherwise)
func (b *BookmarkRepository) isBookmarkCreated(bookmarkedNovel models.BookmarkedNovel) bool {
	var existingBookmarkedNovel models.BookmarkedNovel
	if err := b.db.Where("novel_id = ? AND user_id = ? AND deleted_at IS NOT NULL", bookmarkedNovel.NovelID, bookmarkedNovel.UserID).First(&existingBookmarkedNovel).Error; err != nil {
		return false
	}
	return existingBookmarkedNovel.ID != 0
}

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
func (b *BookmarkRepository) DeleteBookmark(userID uint, novelID uint) error {
	err := b.db.Model(&models.BookmarkedNovel{}).
		Where("user_id = ? AND novel_id = ?", userID, novelID).
		Delete(&models.BookmarkedNovel{}).Error

	if err != nil {
		if err.Error() == "record not found" {
			return types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to delete bookmarked novel", err)
	}
	return nil
}
