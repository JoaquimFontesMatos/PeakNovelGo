package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// NovelRepository struct represents a repository for novel management.
type NovelRepository struct {
	db *gorm.DB
}

// NewNovelRepository creates a new NovelRepository instance
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *NovelRepository (pointer to the NovelRepository instance)
func NewNovelRepository(db *gorm.DB) *NovelRepository {
	return &NovelRepository{db: db}
}

// CreateNovel creates a new novel in the database.
//
// Parameters:
//   - novel models.Novel (Novel struct)
//
// Returns:
//   - *models.Novel (pointer to Novel struct)
//   - CONFLICT_ERROR if the novel already exists
//   - INTERNAL_SERVER_ERROR if the novel could not be created
func (n *NovelRepository) CreateNovel(novel models.Novel) (*models.Novel, error) {
	n.db.Logger = n.db.Logger.LogMode(logger.Silent)

	if IsNovelCreated := n.isNovelCreated(novel); IsNovelCreated {
		return nil, types.WrapError("CONFLICT_ERROR", "Novel already exists", nil)
	}

	// Initialize slices for the new relationships
	var newTags []models.Tag
	var newAuthors []models.Author
	var newGenres []models.Genre

	// Process tags
	for _, tag := range novel.Tags {
		var existingTag models.Tag

		if tag.Name == "" {
			continue // Skip empty tags
		}

		err := n.db.Where("name = ?", tag.Name).FirstOrCreate(&existingTag, models.Tag{
			Name:        tag.Name,
			Description: tag.Description,
		}).Error

		if err != nil {
			return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to associate tag", err)
		}

		newTags = append(newTags, existingTag) // Append to the newTags slice
	}

	// Process authors
	for _, author := range novel.Authors {
		var existingAuthor models.Author

		if author.Name == "" {
			continue // Skip empty authors
		}

		err := n.db.Where("name = ?", author.Name).FirstOrCreate(&existingAuthor, models.Author{
			Name: author.Name,
		}).Error

		if err != nil {
			return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to associate author", err)
		}

		newAuthors = append(newAuthors, existingAuthor) // Append to the newAuthors slice
	}

	// Process genres
	for _, genre := range novel.Genres {
		var existingGenre models.Genre

		if genre.Name == "" {
			continue // Skip empty genres
		}

		err := n.db.Where("name = ?", genre.Name).FirstOrCreate(&existingGenre, models.Genre{
			Name:        genre.Name,
			Description: genre.Description,
		}).Error

		if err != nil {
			return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to associate genre", err)
		}

		newGenres = append(newGenres, existingGenre) // Append to the newGenres slice
	}

	// Update the novel with associated relationships
	novel.Tags = newTags
	novel.Authors = newAuthors
	novel.Genres = newGenres

	// Save the novel with relationships
	if err := n.db.Create(&novel).Error; err != nil {
		log.Println(err)
		return nil, types.WrapError("INTERNAL_SERVER_ERROR", "Failed to create novel", err)
	}

	return &novel, nil
}

// isNovelCreated checks if a novel with the given URL already exists in the database.
//
// Parameters:
//   - novel models.Novel (Novel struct)
//
// Returns:
//   - bool (true if the novel already exists, false otherwise)
func (n *NovelRepository) isNovelCreated(novel models.Novel) bool {
	var existingNovel models.Novel
	if err := n.db.Where("novel_updates_id = ?", novel.NovelUpdatesID).First(&existingNovel).Error; err != nil {
		return false
	}
	return existingNovel.ID != 0
}

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
func (n *NovelRepository) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of novels", err)
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
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novels", err)
	}
	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the genre
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_genres ON novel_genres.novel_id = novels.id").
		Joins("JOIN genres ON genres.id = novel_genres.genre_id").
		Where("genres.name = ?", genreName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of novels", err)
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
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novels", err)
	}

	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the tag
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_tags ON novel_tags.novel_id = novels.id").
		Joins("JOIN tags ON tags.id = novel_tags.tag_id").
		Where("tags.name = ?", tagName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of novels", err)
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
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novels", err)
	}

	return novels, total, nil
}

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
func (n *NovelRepository) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_authors ON novel_authors.novel_id = novels.id").
		Joins("JOIN authors ON authors.id = novel_authors.author_id").
		Where("authors.name = ?", authorName).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of novels", err)
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
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novels", err)
	}
	return novels, total, nil
}

// GetNovelByID gets a novel by ID.
//
// Parameters:
//   - id uint (ID of the novel)
//
// Returns:
//   - *models.Novel (pointer to Novel struct)
//   - INTERNAL_SERVER_ERROR if the novel could not be fetched
//   - NOVEL_NOT_FOUND_ERROR if the novel could not be fetched
func (n *NovelRepository) GetNovelByID(id uint) (*models.Novel, error) {
	var novel models.Novel
	if err := n.db.Where("id = ?", id).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		First(&novel).
		Error; err != nil {

		if err.Error() == "record not found" {
			return nil, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "No novels found", nil)
		}

		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novel", err)
	}
	return &novel, nil
}

// GetNovelByUpdatesID gets a novel by NovelUpdatesID.
//
// Parameters:
//   - NovelUpdatesID string (NovelUpdatesID of the novel)
//
// Returns:
//   - *models.Novel (pointer to Novel struct)
//   - INTERNAL_SERVER_ERROR if the novel could not be fetched
//   - NOVEL_NOT_FOUND_ERROR if the novel could not be fetched
func (n *NovelRepository) GetNovelByUpdatesID(title string) (*models.Novel, error) {
	var novel models.Novel
	if err := n.db.Where("novel_updates_id = ?", title).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		First(&novel).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "No novels found", nil)
		}

		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch novel", err)
	}
	return &novel, nil
}
