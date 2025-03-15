package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"backend/internal/types/errors"
	"log"
	"net/http"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ChapterRepository struct represents a repository for chapter management.
type ChapterRepository struct {
	db *gorm.DB
}

// NewChapterRepository creates a new ChapterRepository instance
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *ChapterRepository (pointer to the ChapterRepository instance)
func NewChapterRepository(db *gorm.DB) *ChapterRepository {
	return &ChapterRepository{db: db}
}

// isChapterCreated checks if a chapter with the given chapter number and novel ID already exists in the database.
//
// Parameters:
//   - chapter models.Chapter (Chapter struct)
//
// Returns:
//   - bool (true if the chapter already exists, false otherwise)
func (c *ChapterRepository) IsChapterCreated(chapterNo uint, novelID uint) bool {
	var existingChapter models.Chapter
	if err := c.db.Where("chapter_no = ? AND novel_id = ?", chapterNo, novelID).First(&existingChapter).Error; err != nil {
		return false
	}
	return existingChapter.ID != 0
}

// CreateChapter creates a new chapter in the database.
//
// Parameters:
//   - chapter models.Chapter (Chapter struct)
//
// Returns:
//   - *models.Chapter (pointer to Chapter struct)
//   - CONFLICT_ERROR if the chapter already exists
//   - INTERNAL_SERVER_ERROR if the chapter could not be created
func (c *ChapterRepository) CreateChapter(chapter models.Chapter) (*models.Chapter, error) {
	c.db.Logger = c.db.Logger.LogMode(logger.Silent)

	if IsChapterCreated := c.IsChapterCreated(chapter.ChapterNo, *chapter.NovelID); IsChapterCreated {
		return nil, errors.ErrChapterConflict
	}

	// Save the chapter
	if err := c.db.Create(&chapter).Error; err != nil {
		log.Println(err)
		return nil, types.WrapError(errors.IMPORTING_CHAPTER, "Failed to create chapter", http.StatusInternalServerError, err)
	}

	return &chapter, nil
}

// GetChapterByNovelUpdatesIDAndChapterNo gets a chapter by novel title and chapter number.
//
// Parameters:
//   - novelTitle string (title of the novel)
//   - chapterNo uint (chapter number)
//
// Returns:
//   - *models.Chapter (pointer to Chapter struct)
//   - INTERNAL_SERVER_ERROR if the chapter could not be fetched
//   - CHAPTER_NOT_FOUND_ERROR if the chapter could not be fetched
func (c *ChapterRepository) GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := c.db.Model(&models.Chapter{}).
		Joins("JOIN novels ON novels.id = chapters.novel_id").
		Where("novels.novel_updates_id = ? AND chapters.chapter_no = ?", novelTitle, chapterNo).First(&chapter).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errors.ErrChapterNotFound
		}
		log.Println(err)
		return nil, types.WrapError(errors.GETTING_CHAPTER, "Failed to fetch chapter", http.StatusInternalServerError, err)
	}
	return &chapter, nil
}

// GetChaptersByNovelUpdatesID gets a list of chapters by novel updates ID.
//
// Parameters:
//   - novelTitle string (title of the novel)
//
// Returns:
//   - []models.Chapter (list of Chapter structs)
//   - INTERNAL_SERVER_ERROR if the chapters could not be fetched
//   - NO_CHAPTERS_ERROR if the chapters could not be fetched
func (c *ChapterRepository) GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error) {
	var chapters []models.Chapter
	var total int64

	if err := c.db.Model(&models.Chapter{}).
		Joins("JOIN novels ON novels.id = chapters.novel_id").
		Where("novels.novel_updates_id = ?", novelTitle).
		Count(&total).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoChapters
		}

		return nil, 0, types.WrapError(errors.GETTING_TOTAL_CHAPTERS, "Failed to get the total number of chapters", http.StatusInternalServerError, err)
	}
	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := c.db.Model(&models.Chapter{}).
		Joins("JOIN novels ON novels.id = chapters.novel_id").
		Where("novels.novel_updates_id = ?", novelTitle).
		Order("chapter_no ASC").
		Limit(limit).Offset(offset).
		Offset(offset).
		Find(&chapters).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoChapters
		}

		return nil, 0, types.WrapError(errors.GETTING_CHAPTERS, "Failed to fetch chapters", http.StatusInternalServerError, err)
	}
	return chapters, total, nil
}
