package interfaces

import (
	"backend/internal/models"
)

type NovelRepositoryInterface interface {
	// IsDown checks if the database is offline by pinging it
	IsDown() bool

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
}
