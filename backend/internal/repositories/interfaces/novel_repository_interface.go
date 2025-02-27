package interfaces

import (
	"backend/internal/models"
)

type NovelRepositoryInterface interface {
	// CreateNovel creates a new novel in the database.
	//
	// Parameters:
	//   - novel models.Novel (Novel struct)
	//
	// Returns:
	//   - *models.Novel (pointer to Novel struct)
	//   - CONFLICT_ERROR if the novel already exists
	//   - INTERNAL_SERVER_ERROR if the novel could not be created
	CreateNovel(novel models.Novel) (*models.Novel, error)

	// CreateChapters creates a list of chapters in the database.
	//
	// Parameters:
	//   - chapters []models.Chapter (list of Chapter structs)
	//
	// Returns:
	//   - int (number of chapters created)
	//   - INTERNAL_SERVER_ERROR if the chapters could not be created
	//   - NO_NEW_CHAPTERS_ERROR if there's no new chapters to create
	CreateChapters(chapters []models.Chapter) (int, error)

	// GetChaptersByNovelID gets a list of chapters by novel ID.
	//
	// Parameters:
	//   - novelID uint (ID of the novel)
	//   - page int (page number)
	//   - limit int (limit of chapters per page)
	//
	// Returns:
	//   - []models.Chapter (list of Chapter structs)
	//   - int64 (total number of chapters)
	//   - INTERNAL_SERVER_ERROR if the chapters could not be fetched
	//   - NO_CHAPTERS_ERROR if the chapters could not be fetched
	GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error)

	// GetNovelsByAuthorName gets a list of novels by author name.
	//
	// Parameters:
	//   - authorName string (name of the author)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovels gets a list of novels.
	//
	// Parameters:
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovels(page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByGenreName gets a list of novels by genre name.
	//
	// Parameters:
	//   - genreName string (name of the genre)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByTagName gets a list of novels by tag name.
	//
	// Parameters:
	//   - tagName string (name of the tag)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelByID gets a novel by ID.
	//
	// Parameters:
	//   - id uint (ID of the novel)
	//
	// Returns:
	//   - *models.Novel (pointer to Novel struct)
	//   - INTERNAL_SERVER_ERROR if the novel could not be fetched
	//   - NOVEL_NOT_FOUND_ERROR if the novel could not be fetched
	GetNovelByID(id uint) (*models.Novel, error)

	// GetNovelByUpdatesID gets a novel by novel updates id.
	//
	// Parameters:
	//   - title string (title of the novel)
	//
	// Returns:
	//   - *models.Novel (pointer to Novel struct)
	//   - INTERNAL_SERVER_ERROR if the novel could not be fetched
	//   - NOVEL_NOT_FOUND_ERROR if the novel could not be fetched
	GetNovelByUpdatesID(title string) (*models.Novel, error)

	// GetChapterByID gets a chapter by ID.
	//
	// Parameters:
	//   - id uint (ID of the chapter)
	//
	// Returns:
	//   - *models.Chapter (pointer to Chapter struct)
	//   - INTERNAL_SERVER_ERROR if the chapter could not be fetched
	//   - CHAPTER_NOT_FOUND_ERROR if the chapter could not be fetched
	GetChapterByID(id uint) (*models.Chapter, error)

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
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error)

	// GetBookmarkedNovelByUserIDAndNovelID gets a bookmarked novel by user ID and novel ID.
	//
	// Parameters:
	//   - userID uint (ID of the user)
	//   - novelID string (ID of the novel)
	//
	// Returns:
	//   - models.BookmarkedNovel (BookmarkedNovel struct)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be fetched
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be fetched
	GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error)

	// UpdateBookmarkedNovel updates a bookmarked novel in the database.
	//
	// Parameters:
	//   - novel models.BookmarkedNovel (BookmarkedNovel struct)
	//
	// Returns:
	//   - models.BookmarkedNovel (BookmarkedNovel struct)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be updated
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be updated
	UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error)

	// CreateBookmarkedNovel creates a new bookmarked novel in the database.
	//
	// Parameters:
	//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
	//
	// Returns:
	//   - *models.BookmarkedNovel (pointer to BookmarkedNovel struct)
	//   - CONFLICT_ERROR if the bookmarked novel already exists
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be created
	CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)

	// DeleteBookmarkedNovel deletes a bookmarked novel from the database.
	//
	// Parameters:
	//   - userID uint (ID of the user)
	//   - novelID uint (ID of the novel)
	//
	// Returns:
	//   - error (nil if the bookmarked novel was deleted successfully, otherwise an error)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be deleted
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be deleted
	DeleteBookmarkedNovel(userID uint, novelID uint) error
}
