package interfaces

import "backend/internal/models"

type ChapterRepositoryInterface interface {
	// isChapterCreated checks if a chapter with the given chapter number and novel ID already exists in the database.
	//
	// Parameters:
	//   - chapterNo uint (chapter number)
	//   - novelID uint (ID of the novel)
	//
	// Returns:
	//   - bool (true if the chapter already exists, false otherwise)
	//   - INTERNAL_SERVER_ERROR if the chapter could not be fetched
	IsChapterCreated(chapterNo uint, novelID uint) bool

	// CreateChapter creates a new chapter in the database.
	//
	// Parameters:
	//   - chapter models.Chapter (Chapter struct)
	//
	// Returns:
	//   - *models.Chapter (pointer to Chapter struct)
	//   - CONFLICT_ERROR if the chapter already exists
	//   - INTERNAL_SERVER_ERROR if the chapter could not be created
	CreateChapter(chapter models.Chapter) (*models.Chapter, error)

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
	GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error)

	// GetChaptersByNovelUpdatesID gets a list of chapters by novel title and chapter number.
	//
	// Parameters:
	//   - novelTitle string (title of the novel)
	//   - chapterNo uint (chapter number)
	//
	// Returns:
	//   - []models.Chapter (list of Chapter structs)
	//   - INTERNAL_SERVER_ERROR if the chapters could not be fetched
	//   - NO_CHAPTERS_ERROR if the chapters could not be fetched
	GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error)
}
