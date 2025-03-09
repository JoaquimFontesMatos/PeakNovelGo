package interfaces

import (
	"backend/internal/models"
)

// NovelServiceInterface defines methods for managing and retrieving novel-related data within the application.
type NovelServiceInterface interface {
	// GetNovelsByAuthorName retrieves novels by author name, handling pagination.
	//
	// Parameters:
	//   - authorName (string): The name of the author to search for.  Must be a valid author name.
	//   - page (int): The page number for pagination (starting from 1).
	//   - limit (int): The number of novels to retrieve per page.
	//
	// Returns:
	//   - []models.Novel: A slice of Novel structs matching the criteria.  An empty slice if no novels are found.
	//   - int64: The total number of novels matching the criteria (across all pages).
	//   - error: An error object if something goes wrong (e.g., invalid author name, database error).  Nil on success.
	//
	// Error types:
	//
	// Validation errors:
	//   - errors.ErrAuthorRequired: Returned if the author string is empty or contains only spaces.
	//   - errors.ErrAuthorTooShort: Returned if the author string is less than 1 character long.
	//   - errors.ErrAuthorTooLong: Returned if the author string is more than 255 characters long.
	//
	// Repository errors:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given author.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByGenreName retrieves novels by genre name, handling pagination.
	//
	// Parameters:
	//   - genreName (string): The name of the genre to filter by.  Must be a valid genre name.
	//   - page (int): The page number for pagination (starting from 1).
	//   - limit (int): The number of novels per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels matching the criteria.  An empty slice is returned if no novels are found.
	//   - int64: The total number of novels matching the criteria (regardless of pagination).
	//   - error: An error if the genre name is invalid or a database error occurs.  Returns a validation error if genreName is invalid.
	//
	// Error types:
	//
	// Validation errors:
	//   - errors.ErrGenreRequired: if the genre string is empty or contains only spaces.
	//   - errors.ErrGenreTooShort: if the genre string is less than 1 character long.
	//   - errors.ErrGenreTooLong: if the genre string is more than 255 characters long.
	//
	// Repository errors:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given author.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovelsByTagName retrieves novels by tag name, handling pagination.
	//
	// Parameters:
	//   - tagName (string): The name of the tag to filter by.  Must be a valid tag name.
	//   - page (int): The page number for pagination (starting from 1).
	//   - limit (int): The number of novels per page.
	//
	// Returns:
	//   - []models.Novel: A slice of novels matching the criteria.  An empty slice is returned if no novels are found.
	//   - int64: The total number of novels matching the criteria (regardless of pagination).
	//   - error: An error if the tag name is invalid or a database error occurs.  Returns a validation error if tagName is invalid.
	//
	// Error types:
	//
	// Validation errors:
	//   - errors.ErrTagRequired: if the tag string is empty or contains only spaces.
	//   - errors.ErrTagTooShort: if the tag string is less than 1 character long.
	//   - errors.ErrTagTooLong: if the tag string is more than 255 characters long.
	//
	// Repository errors:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found for the given author.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error)

	// GetNovels retrieves a paginated list of novels.
	//
	// Parameters:
	//   - page (int): The page number (starting from 1).
	//   - limit (int): The number of novels per page.
	//
	// Returns:
	//   - []models.Novel: A slice of Novel structs representing the novels on the specified page.
	//   - int64: The total number of novels.
	//   - error: An error object; returns nil if no error occurs.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	GetNovels(page, limit int) ([]models.Novel, int64, error)

	// GetNovelByID retrieves a novel from the repository based on its ID.
	//
	// Parameters:
	//   - id (uint): The ID of the novel to retrieve.
	//
	// Returns:
	//   - *models.Novel: A pointer to the retrieved novel, or nil if not found.
	//   - error: An error object if any issues occur during retrieval.  May indicate database errors or other unexpected issues.
	//
	// Error types:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNovelNotFound: Returned if the novel with the given ID is not found.
	//   - errors.ErrGettingNovel: Returned if there's a general error retrieving the novel.
	GetNovelByID(id uint) (*models.Novel, error)

	// GetNovelByUpdatesID retrieves a novel from the repository using its NovelUpdates ID.
	//
	// Parameters:
	//   - novelUpdatesID (string): The NovelUpdates ID of the novel to retrieve.
	//
	// Returns:
	//   - *models.Novel: A pointer to the retrieved novel, or nil if not found.
	//   - error: An error if the NovelUpdates ID is invalid or if there's an issue retrieving the novel.  The error will
	//     indicate the specific problem encountered.
	//
	// Error types:
	//
	// Validation errors:
	//   - errors.ErrInvalidNovelUpdatesID: Returned if the provided ID is invalid (e.g., incorrect length or format).
	//
	// Repository errors:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNovelNotFound: Returned if no novel with the given NovelUpdates ID exists in the database.
	//   - errors.ErrGettingNovel: Returned if an error occurred while retrieving the novel from the database.
	GetNovelByUpdatesID(title string) (*models.Novel, error)

	// CreateNovel creates a new novel in the database using data scraped from NovelUpdates.
	//
	// It takes a NovelUpdates ID, parses it, executes a Python script to scrape novel data,
	// handles potential errors during script execution and JSON unmarshalling,
	// cleans the scraped data, and finally saves the novel to the database.
	//
	// Parameters:
	//   - novelUpdatesID (string): The ID of the novel on NovelUpdates.
	//
	// Returns:
	//   - *models.Novel: A pointer to the newly created novel.  Returns nil if an error occurs.
	//   - error: An error object if any error occurred during the process.  Specific error types are detailed below.
	//
	// Error types:
	//
	// Service errors:
	//   - errors.ErrNovelNotFound: Returned if the novel is not found on NovelUpdates.
	//   - errors.SCRIPT_ERROR:  Indicates failure during Python script execution.  HTTP Status: 503
	//   - errors.IMPORTING_NOVEL: Indicates failure during JSON unmarshalling of the script's output. HTTP Status: 500
	//
	// Validation errors:
	//   - errors.ErrInvalidLatestChapter: Returned if the latest chapter value is invalid.
	//
	// Repository errors:
	//   - errors.ErrDatabaseOffline: Returned if the database is offline.
	//   - errors.ErrNoNovels: Returned if no novels are found.
	//   - errors.ErrGettingTotalNovels: Returned if an error occurs while retrieving the total number of novels.
	//   - errors.ErrGettingNovels: Returned if an error occurs while retrieving the novels themselves.
	CreateNovel(novelUpdatesID string) (*models.Novel, error)
}
