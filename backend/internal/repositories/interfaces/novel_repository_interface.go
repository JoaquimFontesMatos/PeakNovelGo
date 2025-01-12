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

	// ConvertToNovel convert an imported novel struct to a novel struct.
	//
	// Parameters:
	//   - imported models.ImportedNovel (imported novel struct)
	//
	// Returns:
	//   - *models.Novel (pointer to Novel struct)
	//   - INTERNAL_SERVER_ERROR if an error occurred while converting the imported novel struct to a novel struct
	ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error)

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

	// GetNovelsByAuthorID gets a list of novels by author ID.
	//
	// Parameters:
	//   - authorID uint (ID of the author)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error)

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

	// GetNovelsByGenreID gets a list of novels by genre ID.
	//
	// Parameters:
	//   - genreID uint (ID of the genre)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByTagID gets a list of novels by tag ID.
	//
	// Parameters:
	//   - tagID uint (ID of the tag)
	//   - page int (page number)
	//   - limit int (limit of novels per page)
	//
	// Returns:
	//   - []models.Novel (list of Novel structs)
	//   - int64 (total number of novels)
	//   - INTERNAL_SERVER_ERROR if the novels could not be fetched
	//   - NO_NOVELS_ERROR if the novels could not be fetched
	GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error)

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
	//   - novelID uint (ID of the novel)
	//
	// Returns:
	//   - models.BookmarkedNovel (BookmarkedNovel struct)
	//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be fetched
	//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be fetched
	GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error)

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
