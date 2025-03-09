package interfaces

import (
	"backend/internal/models"
)

// NovelRepositoryInterface defines the contract for managing novel data in the repository layer.
type NovelRepositoryInterface interface {
	BaseRepositoryInterface

	// CreateNovel creates a new novel in the database, handling associated tags, authors, and genres.  It checks for existing
	// records and creates them if necessary, preventing duplicates.
	//
	// Parameters:
	//   - novel (models.Novel): The novel data to be created.  This includes tags, authors, and genres which will be processed.
	//
	// Returns:
	//   - *models.Novel: A pointer to the newly created novel.  Returns nil if an error occurs.
	//   - error: An error object indicating the failure, or nil if successful.  Specific error types are detailed below.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Indicates that the database is offline and cannot be accessed.
	//   - errors.ErrNovelAlreadyImported: Indicates that a novel with the same data already exists.
	//   - types.WrappedError (errors.TAG_ASSOCIATION_ERROR): Wrapped error indicating failure during tag creation or association.
	//   - types.WrappedError (errors.AUTHOR_ASSOCIATION_ERROR): Wrapped error indicating failure during author creation or association.
	//   - types.WrappedError (errors.GENRE_ASSOCIATION_ERROR): Wrapped error indicating failure during genre creation or association.
	//   - errors.ErrImportingNovel: A generic error indicating failure during novel creation.
	CreateNovel(novel models.Novel) (*models.Novel, error)

	// GetNovels retrieves a paginated list of all the novels
	//
	// Parameters:
	//   - page (int): The page number for pagination (1-based).
	//   - limit (int): The maximum number of novels to return per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels.
	//   - int64: The total number of novels.
	//   - error: An error object indicating any issues encountered during the retrieval process.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovels(page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByAuthorName retrieves a paginated list of novels by a given author name.
	//
	// Parameters:
	//   - authorName (string): The name of the author to search for.
	//   - page (int): The page number for pagination (1-based).
	//   - limit (int): The maximum number of novels to return per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels matching the author's name.
	//   - int64: The total number of novels by the author.
	//   - error: An error object indicating any issues encountered during the retrieval process.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given author.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByGenreName retrieves a paginated list of novels by a given genre.
	//
	// Parameters:
	//   - genreName (string): The name of the genre to search for.
	//   - page (int): The page number for pagination (1-based).
	//   - limit (int): The maximum number of novels to return per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels matching the genre.
	//   - int64: The total number of novels by the genre.
	//   - error: An error object indicating any issues encountered during the retrieval process.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given genre.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByTagName retrieves a paginated list of novels by a given tag.
	//
	// Parameters:
	//   - tagName (string): The name of the tag to search for.
	//   - page (int): The page number for pagination (1-based).
	//   - limit (int): The maximum number of novels to return per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels matching the tag.
	//   - int64: The total number of novels by the tag.
	//   - error: An error object indicating any issues encountered during the retrieval process.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given tag.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelByID retrieves a novel from the database based on its ID.
	// It preloads the associated authors, genres, and tags.
	//
	// Parameters:
	//   - id (uint): The ID of the novel to retrieve.
	//
	// Returns:
	//   - *models.Novel: A pointer to the retrieved novel, or nil if not found.
	//   - error: An error object indicating the reason for failure, if any.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNovelNotFound: Returned if the novel with the given ID is not found.
	//   - errors.ErrGettingNovel: Returned if there's a general error retrieving the novel.
	GetNovelByID(id uint) (*models.Novel, error)

	// GetNovelByUpdatesID retrieves a novel from the database based on its NovelUpdates ID.
	//
	// Parameters:
	//   - title (string): The NovelUpdates ID of the novel.
	//
	// Returns:
	//   - *models.Novel: A pointer to the retrieved novel, or nil if not found.
	//   - error: An error object indicating the type of error encountered, or nil if the operation was successful.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNovelNotFound: Returned if no novel with the given NovelUpdates ID exists in the database.
	//   - errors.ErrGettingNovel: Returned if an error occurred while retrieving the novel from the database.
	GetNovelByUpdatesID(title string) (*models.Novel, error)
}
