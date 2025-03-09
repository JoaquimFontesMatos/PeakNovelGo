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

// NovelRepository represents a repository for interacting with novel data.
// It embeds the BaseRepository to inherit common database operations.
type NovelRepository struct {
	*BaseRepository
}

// NewNovelRepository creates a new NovelRepository.
//
// Parameters:
//   - db (*gorm.DB): The database connection.
//
// Returns:
//   - *NovelRepository: A pointer to the newly created NovelRepository.
func NewNovelRepository(db *gorm.DB) *NovelRepository {
	return &NovelRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

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
func (n *NovelRepository) CreateNovel(novel models.Novel) (*models.Novel, error) {
	if n.IsDown() {
		return nil, errors.ErrDatabaseOffline
	}

	n.db.Logger = n.db.Logger.LogMode(logger.Silent)

	if IsNovelCreated := n.isNovelCreated(novel); IsNovelCreated {
		return nil, errors.ErrNovelAlreadyImported
	}

	// Process tags
	newTags, err := n.processTags(novel.Tags)
	if err != nil {
		return nil, err
	}

	// Process authors
	newAuthors, err := n.processAuthors(novel.Authors)
	if err != nil {
		return nil, err
	}

	// Process genres
	newGenres, err := n.processGenres(novel.Genres)
	if err != nil {
		return nil, err
	}

	// Update the novel with associated relationships
	novel.Tags = newTags
	novel.Authors = newAuthors
	novel.Genres = newGenres

	// Save the novel with relationships
	if err := n.db.Create(&novel).Error; err != nil {
		log.Println(err)
		return nil, errors.ErrImportingNovel
	}

	return &novel, nil
}

// isNovelCreated checks if a novel with the given URL already exists in the database.
//
// Parameters:
//   - novel (models.Novel): Novel struct
//
// Returns:
//   - bool: true if the novel already exists, false otherwise
func (n *NovelRepository) isNovelCreated(novel models.Novel) bool {
	var existingNovel models.Novel
	if err := n.db.Where("novel_updates_id = ?", novel.NovelUpdatesID).First(&existingNovel).Error; err != nil {
		return false
	}
	return existingNovel.ID != 0
}

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
func (n *NovelRepository) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	if n.IsDown() {
		return nil, 0, errors.ErrDatabaseOffline
	}

	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingTotalNovels
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.Novel{}).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		Limit(limit).Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingNovels
	}
	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	if n.IsDown() {
		return nil, 0, errors.ErrDatabaseOffline
	}

	var novels []models.Novel
	var total int64

	// Count total novels for the genre
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_genres ON novel_genres.novel_id = novels.id").
		Joins("JOIN genres ON genres.id = novel_genres.genre_id").
		Where("genres.name = ?", genreName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingTotalNovels
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_genres ON novel_genres.novel_id = novels.id").
		Joins("JOIN genres ON genres.id = novel_genres.genre_id").
		Where("genres.name = ?", genreName).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingNovels
	}

	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	if n.IsDown() {
		return nil, 0, errors.ErrDatabaseOffline
	}

	var novels []models.Novel
	var total int64

	// Count total novels for the tag
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_tags ON novel_tags.novel_id = novels.id").
		Joins("JOIN tags ON tags.id = novel_tags.tag_id").
		Where("tags.name = ?", tagName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingTotalNovels
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_tags ON novel_tags.novel_id = novels.id").
		Joins("JOIN tags ON tags.id = novel_tags.tag_id").
		Where("tags.name = ?", tagName).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		Limit(limit).
		Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingNovels
	}

	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	if n.IsDown() {
		return nil, 0, errors.ErrDatabaseOffline
	}

	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_authors ON novel_authors.novel_id = novels.id").
		Joins("JOIN authors ON authors.id = novel_authors.author_id").
		Where("authors.name = ?", authorName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingTotalNovels
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_authors ON novel_authors.novel_id = novels.id").
		Joins("JOIN authors ON authors.id = novel_authors.author_id").
		Where("authors.name = ?", authorName).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		Limit(limit).Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, errors.ErrNoNovels
		}

		return nil, 0, errors.ErrGettingNovels
	}
	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelByID(id uint) (*models.Novel, error) {
	if n.IsDown() {
		return nil, errors.ErrDatabaseOffline
	}

	var novel models.Novel
	if err := n.db.Where("id = ?", id).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		First(&novel).
		Error; err != nil {

		if err.Error() == "record not found" {
			return nil, errors.ErrNovelNotFound
		}

		return nil, errors.ErrGettingNovel
	}
	return &novel, nil
}

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
func (n *NovelRepository) GetNovelByUpdatesID(title string) (*models.Novel, error) {
	if n.IsDown() {
		return nil, errors.ErrDatabaseOffline
	}

	var novel models.Novel
	if err := n.db.Where("novel_updates_id = ?", title).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		First(&novel).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, errors.ErrNovelNotFound
		}

		return nil, errors.ErrGettingNovel
	}
	return &novel, nil
}

// processTags processes a list of tags, ensuring they exist in the database.
// It iterates through the input tags and either retrieves existing tags or creates new ones if they don't exist.
// Empty tag names are skipped.
//
// Parameters:
//   - tags ([]models.Tag): A slice of tag structs to process.
//
// Returns:
//   - []models.Tag: A slice of processed tags, containing either existing or newly created tags.  Each tag will
//     have a valid ID.
//   - error: An error if there's a problem creating or retrieving a tag from the database.  The error will wrap a more
//     specific error for better diagnostics.
//
// Error types:
//   - types.Error: Wraps errors.TAG_ASSOCIATION_ERROR with a descriptive message and HTTP status code
//     (http.StatusInternalServerError) if there's an issue with database operations related to tag creation or retrieval.
func (n *NovelRepository) processTags(tags []models.Tag) ([]models.Tag, error) {
	var newTags []models.Tag

	for _, tag := range tags {
		if tag.Name == "" {
			continue // Skip empty tags
		}

		var existingTag models.Tag
		err := n.db.Where("name = ?", tag.Name).FirstOrCreate(&existingTag, models.Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}).Error

		if err != nil {
			return nil, types.WrapError(errors.TAG_ASSOCIATION_ERROR, "Failed to create tag", http.StatusInternalServerError, err)
		}

		newTags = append(newTags, existingTag)
	}

	return newTags, nil
}

// processAuthors processes a list of authors, ensuring they exist in the database.
// It iterates through the input authors and either retrieves existing authors or creates new ones if they don't exist.
// Empty author names are skipped.
//
// Parameters:
//   - authors ([]models.Author): A slice of authors structs to process.
//
// Returns:
//   - []models.Author: A slice of processed authors, containing either existing or newly created authors.  Each author will
//     have a valid ID.
//   - error: An error if there's a problem creating or retrieving an author from the database.  The error will wrap a more
//     specific error for better diagnostics.
//
// Error types:
//   - types.Error: Wraps errors.AUTHOR_ASSOCIATION_ERROR with a descriptive message and HTTP status code
//     (http.StatusInternalServerError) if there's an issue with database operations related to author creation or retrieval.func (n *NovelRepository) processAuthors(authors []models.Author) ([]models.Author, error) {
func (n *NovelRepository) processAuthors(authors []models.Author) ([]models.Author, error) {
	var newAuthors []models.Author

	for _, author := range authors {
		if author.Name == "" {
			continue // Skip empty authors
		}

		var existingAuthor models.Author
		err := n.db.Where("name = ?", author.Name).FirstOrCreate(&existingAuthor, models.Author{
			Name: author.Name,
		}).Error

		if err != nil {
			return nil, types.WrapError(errors.AUTHOR_ASSOCIATION_ERROR, "Failed to create author", http.StatusInternalServerError, err)
		}

		newAuthors = append(newAuthors, existingAuthor)
	}

	return newAuthors, nil
}

// processGenres processes a list of genres, ensuring they exist in the database.
// It iterates through the input genres and either retrieves existing genres or creates new ones if they don't exist.
// Empty genre names are skipped.
//
// Parameters:
//   - genres ([]models.Genre): A slice of genre structs to process.
//
// Returns:
//   - []models.Genre: A slice of processed genres, containing either existing or newly created genres.  Each genre will
//     have a valid ID.
//   - error: An error if there's a problem creating or retrieving a genre from the database.  The error will wrap a more
//     specific error for better diagnostics.
//
// Error types:
//   - types.Error: Wraps errors.GENRE_ASSOCIATION_ERROR with a descriptive message and HTTP status code
//     (http.StatusInternalServerError) if there's an issue with database operations related to genre creation or retrieval.
func (n *NovelRepository) processGenres(genres []models.Genre) ([]models.Genre, error) {
	var newGenres []models.Genre

	for _, genre := range genres {
		if genre.Name == "" {
			continue // Skip empty genres
		}

		var existingGenre models.Genre
		err := n.db.Where("name = ?", genre.Name).FirstOrCreate(&existingGenre, models.Genre{
			Name:        genre.Name,
			Description: genre.Description,
		}).Error

		if err != nil {
			return nil, types.WrapError(errors.GENRE_ASSOCIATION_ERROR, "Failed to create genre", http.StatusInternalServerError, err)
		}

		newGenres = append(newGenres, existingGenre)
	}

	return newGenres, nil
}
