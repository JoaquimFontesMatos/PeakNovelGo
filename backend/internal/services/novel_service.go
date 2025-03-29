package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/types/errors"
	"backend/internal/utils"
	"backend/internal/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// NovelService struct manages novel-related operations.
type NovelService struct {
	repo           interfaces.NovelRepositoryInterface
	scriptExecutor utils.ScriptExecutor
}

// NewNovelService creates a new NovelService instance.
//
// Parameters:
//   - repo (interfaces.NovelRepositoryInterface): The novel repository to use.
//   - scriptExecutor (utils.ScriptExecutor): The script executor to use.
//
// Returns:
//   - *NovelService: A new NovelService instance.
func NewNovelService(repo interfaces.NovelRepositoryInterface, scriptExecutor utils.ScriptExecutor) *NovelService {
	return &NovelService{repo: repo, scriptExecutor: scriptExecutor}
}

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
func (s *NovelService) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateAuthor(authorName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByAuthorName(authorName, page, limit)
}

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
func (s *NovelService) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateGenre(genreName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByGenreName(genreName, page, limit)
}

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
func (s *NovelService) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateTag(tagName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByTagName(tagName, page, limit)
}

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
func (s *NovelService) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovels(page, limit)
}

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
func (s *NovelService) GetNovelByUpdatesID(novelUpdatesID string) (*models.Novel, error) {
	parsedNovelUpdatesID, err := utils.NewNovelUpdatesIDParser().Parse(novelUpdatesID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetNovelByUpdatesID(parsedNovelUpdatesID)
}

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
func (s *NovelService) GetNovelByID(id uint) (*models.Novel, error) {
	return s.repo.GetNovelByID(id)
}

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
func (s *NovelService) CreateNovel(novelUpdatesID string) (*models.Novel, error) {
	novelUpdatesID, err := utils.NewNovelUpdatesIDParser().Parse(novelUpdatesID)
	if err != nil {
		return nil, err
	}

	// Execute the Python script
	output, err := s.scriptExecutor.ExecuteScript(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-novel", novelUpdatesID)

	if err != nil {
		return nil, types.WrapError(errors.SCRIPT_ERROR, "Failed to execute Python script: "+err.Error(), http.StatusServiceUnavailable, err)
	}

	var scriptError utils.ScriptError

	// Check if the script returned a specific error
	if json.Unmarshal(output, &scriptError) == nil {
		if scriptError.Status == 404 {
			return nil, errors.ErrNovelNotFound
		}
		return nil, types.WrapError(errors.SCRIPT_ERROR, scriptError.Error, http.StatusServiceUnavailable, nil)
	}

	// Ensure the output is unmarshaled into a valid JSON object
	var result models.ImportedNovel
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, types.WrapError(errors.IMPORTING_NOVEL, "An error occurred while importing the novel: "+err.Error(), http.StatusInternalServerError, err)
	}

	year := strings.ReplaceAll(result.Year, "\n", "")
	status := strings.ReplaceAll(result.Status, "\n", "")
	language := strings.ReplaceAll(result.Language.Name, "\n", "")

	novel := models.Novel{
		Title:            result.Title,
		Synopsis:         result.Synopsis,
		CoverUrl:         result.CoverUrl,
		Language:         language,
		Status:           status,
		NovelUpdatesUrl:  fmt.Sprintf("https://www.lightnovelworld.co/novel/%s", novelUpdatesID),
		NovelUpdatesID:   novelUpdatesID,
		Tags:             result.Tags,
		Authors:          result.Authors,
		Genres:           result.Genres,
		Year:             year,
		ReleaseFrequency: result.ReleaseFrequency,
		LatestChapter:    result.LatestChapter,
	}

	// Save the novel to the database
	return s.repo.CreateNovel(novel)
}
