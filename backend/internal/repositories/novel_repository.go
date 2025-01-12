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
	if IsNovelCreated := n.isNovelCreated(novel); IsNovelCreated {
		return nil, types.WrapError("CONFLICT_ERROR", "Novel already exists", nil)
	}

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
	if err := n.db.Where("url = ?", novel.Url).First(&existingNovel).Error; err != nil {
		return false
	}
	return existingNovel.ID != 0
}

// getExistingChapterNumbers gets a map of chapter numbers and whether they already exist in the database.
//
// Parameters:
//   - chapterNums []uint (list of chapter numbers)
//   - novelID *uint (ID of the novel)
//
// Returns:
//   - map[string]bool (map of chapter numbers and whether they already exist)
//   - INTERNAL_SERVER_ERROR if the chapter numbers could not be fetched
func (n *NovelRepository) getExistingChapterNumbers(chapterNums []uint, novelID *uint) (map[string]bool, error) {
	var chapters []models.Chapter
	if err := n.db.Select("chapter_no").Where("chapter_no IN ? AND novel_id = ?", chapterNums, novelID).Find(&chapters).Error; err != nil {
		return nil, types.WrapError("INTERNAL_SERVER_ERROR", "Failed to get the existing chapter numbers", err)
	}

	// Create a map for quick lookups
	existingURLs := make(map[string]bool, len(chapters))
	for _, chapter := range chapters {
		existingURLs[chapter.ChapterUrl] = true
	}
	return existingURLs, nil
}

// CreateChapters creates a list of chapters in the database.
//
// Parameters:
//   - chapters []models.Chapter (list of Chapter structs)
//
// Returns:
//   - int (number of chapters created)
//   - INTERNAL_SERVER_ERROR if the chapters could not be created
//   - NO_NEW_CHAPTERS_ERROR if there's no new chapters to create
func (n *NovelRepository) CreateChapters(chapters []models.Chapter) (int, error) {
	length := 0
	n.db.Logger = n.db.Logger.LogMode(logger.Silent)

	return length, n.db.Transaction(func(tx *gorm.DB) error {
		// Filter existing chapters
		chaptersNums := make([]uint, len(chapters))
		for i, chapter := range chapters {
			chaptersNums[i] = chapter.ChapterNo
		}

		if len(chapters) == 0 {
			return types.WrapError("NO_NEW_CHAPTERS_ERROR", "There's no chapters to create", nil)
		}
		novelID := chapters[0].NovelID

		existingURLs, err := n.getExistingChapterNumbers(chaptersNums, novelID)
		if err != nil {
			log.Println(err)
			return err
		}

		// Filter out existing chapters
		newChapters := []models.Chapter{}
		for _, chapter := range chapters {
			if !existingURLs[chapter.ChapterUrl] {
				newChapters = append(newChapters, chapter)
			}
		}

		// Return early if no new chapters to save
		if len(newChapters) == 0 {
			return types.WrapError(types.NO_NEW_CHAPTERS_ERROR, "There's no new chapters to create", nil)
		}

		if err := tx.Create(&newChapters).Error; err != nil {
			return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to create chapters", err)
		}
		length = len(newChapters)

		log.Printf("%d chapters added to the database", length)

		return nil
	})
}

// getTagsByName gets a map of tags and whether they already exist in the database.
//
// Parameters:
//   - names []string (list of tag names)
//
// Returns:
//   - map[string]models.Tag (map of tags and whether they already exist)
//   - INTERNAL_SERVER_ERROR if the tags could not be fetched
func (n *NovelRepository) getTagsByName(names []string) (map[string]models.Tag, error) {
	var existingTags []models.Tag
	if err := n.db.Where("name IN ?", names).Find(&existingTags).Error; err != nil {
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the existing tags", err)
	}

	// Create a map for quick lookups
	tagMap := make(map[string]models.Tag, len(existingTags))
	for _, tag := range existingTags {
		tagMap[tag.Name] = tag
	}
	return tagMap, nil
}

// getAuthorsByName gets a map of authors and whether they already exist in the database.
//
// Parameters:
//   - names []string (list of author names)
//
// Returns:
//   - map[string]models.Author (map of authors and whether they already exist)
//   - INTERNAL_SERVER_ERROR if the authors could not be fetched
func (n *NovelRepository) getAuthorsByName(names []string) (map[string]models.Author, error) {
	var existingAuthors []models.Author
	if err := n.db.Where("name IN ?", names).Find(&existingAuthors).Error; err != nil {
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the existing authors", err)
	}

	// Create a map for quick lookups
	authorMap := make(map[string]models.Author, len(existingAuthors))
	for _, author := range existingAuthors {
		authorMap[author.Name] = author
	}
	return authorMap, nil
}

// getGenresByName gets a map of genres and whether they already exist in the database.
//
// Parameters:
//   - names []string (list of genre names)
//
// Returns:
//   - map[string]models.Genre (map of genres and whether they already exist)
//   - INTERNAL_SERVER_ERROR if the genres could not be fetched
func (n *NovelRepository) getGenresByName(names []string) (map[string]models.Genre, error) {
	var existingGenres []models.Genre
	if err := n.db.Where("name IN ?", names).Find(&existingGenres).Error; err != nil {
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the existing genres", err)
	}

	// Create a map for quick lookups
	genreMap := make(map[string]models.Genre, len(existingGenres))
	for _, genre := range existingGenres {
		genreMap[genre.Name] = genre
	}
	return genreMap, nil
}

// ConvertToNovel convert an imported novel struct to a novel struct.
//
// Parameters:
//   - imported models.ImportedNovel (imported novel struct)
//
// Returns:
//   - *models.Novel (pointer to Novel struct)
//   - INTERNAL_SERVER_ERROR if an error occurred while converting the imported novel struct to a novel struct
func (n *NovelRepository) ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error) {
	// Handle Tags
	tagMap, err := n.getTagsByName(imported.Tags)
	if err != nil {
		return nil, err
	}
	var tags []models.Tag
	for _, tagName := range imported.Tags {
		if tag, exists := tagMap[tagName]; exists {
			tags = append(tags, tag)
		} else {
			tags = append(tags, models.Tag{Name: tagName})
		}
	}

	// Handle Authors
	authorMap, err := n.getAuthorsByName(imported.Authors)
	if err != nil {
		return nil, err
	}
	var authors []models.Author
	for _, authorName := range imported.Authors {
		if author, exists := authorMap[authorName]; exists {
			authors = append(authors, author)
		} else {
			authors = append(authors, models.Author{Name: authorName})
		}
	}

	// Handle Genres
	genreMap, err := n.getGenresByName(imported.Genres)
	if err != nil {
		return nil, err
	}
	var genres []models.Genre
	for _, genreName := range imported.Genres {
		if genre, exists := genreMap[genreName]; exists {
			genres = append(genres, genre)
		} else {
			genres = append(genres, models.Genre{Name: genreName})
		}
	}

	// Return the Novel struct with associated tags, authors, and genres
	return &models.Novel{
		Url:             imported.Url,
		Title:           imported.Title,
		Synopsis:        imported.Synopsis,
		CoverUrl:        imported.CoverUrl,
		Language:        imported.Language,
		Status:          imported.Status,
		NovelUpdatesUrl: imported.NovelUpdatesUrl,
		Tags:            tags,
		Authors:         authors,
		Genres:          genres,
	}, nil
}

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
func (n *NovelRepository) GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error) {
	var chapters []models.Chapter
	var total int64

	// Count total chapters for the novel
	if err := n.db.Model(&models.Chapter{}).Where("novel_id = ?", novelID).Count(&total).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_CHAPTERS_ERROR, "No chapters found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of chapters", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Where("novel_id = ?", novelID).
		Order("chapter_no ASC").
		Limit(limit).
		Offset(offset).
		Find(&chapters).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_CHAPTERS_ERROR, "No chapters found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch chapters", err)
	}

	return chapters, total, nil
}

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
func (n *NovelRepository) GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the genre
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_genres ON novel_genres.novel_id = novels.id").
		Where("novel_genres.genre_id = ?", genreID).
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
		Where("novel_genres.genre_id = ?", genreID).
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
func (n *NovelRepository) GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the tag
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_tags ON novel_tags.novel_id = novels.id").
		Where("novel_tags.tag_id = ?", tagID).
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
		Where("novel_tags.tag_id = ?", tagID).
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
func (n *NovelRepository) GetNovels( page, limit int) ([]models.Novel, int64, error) {
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
func (n *NovelRepository) GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_authors ON novel_authors.novel_id = novels.id").
		Where("novel_authors.author_id = ?", authorID).
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
		Where("novel_authors.author_id = ?", authorID).
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

// GetChapterByID gets a chapter by ID.
//
// Parameters:
//   - id uint (ID of the chapter)
//
// Returns:
//   - *models.Chapter (pointer to Chapter struct)
//   - INTERNAL_SERVER_ERROR if the chapter could not be fetched
//   - CHAPTER_NOT_FOUND_ERROR if the chapter could not be fetched
func (n *NovelRepository) GetChapterByID(id uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := n.db.Where("id = ?", id).First(&chapter).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, types.WrapError(types.CHAPTER_NOT_FOUND_ERROR, "Chapter not found", nil)
		}

		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch chapter", err)
	}
	return &chapter, nil
}

// isChapterCreated checks if a chapter with the given chapter number and novel ID already exists in the database.
//
// Parameters:
//   - chapter models.Chapter (Chapter struct)
//
// Returns:
//   - bool (true if the chapter already exists, false otherwise)
func (n *NovelRepository) isChapterCreated(chapter models.Chapter) bool {
	var existingChapter models.Chapter
	if err := n.db.Where("chapter_no = ? AND novel_id = ?", chapter.ChapterNo, chapter.NovelID).First(&existingChapter).Error; err != nil {
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
func (n *NovelRepository) CreateChapter(chapter models.Chapter) (*models.Chapter, error) {
	if IsChapterCreated := n.isChapterCreated(chapter); IsChapterCreated {
		return nil, types.WrapError(types.CONFLICT_ERROR, "Chapter already exists", nil)
	}

	// Save the chapter
	if err := n.db.Create(&chapter).Error; err != nil {
		log.Println(err)
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to create chapter", err)
	}

	return &chapter, nil
}

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
func (n *NovelRepository) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error) {
	var novels []models.BookmarkedNovel
	var total int64

	// Count total novels for the user
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ?", userID).
		Count(&total).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of bookmarked novels", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&novels).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_NOVELS_ERROR, "No novels found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch bookmarked novels", err)
	}

	return novels, total, nil
}

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
func (n *NovelRepository) GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error) {
	var novel models.BookmarkedNovel
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ? AND bookmarked_novels.novel_id = ?", userID, novelID).
		First(&novel).Error; err != nil {

		if err.Error() == "record not found" {
			return novel, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return novel, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch bookmarked novel", err)
	}
	return novel, nil
}

// UpdateBookmarkedNovel updates a bookmarked novel in the database.
//
// Parameters:
//   - novel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - models.BookmarkedNovel (BookmarkedNovel struct)
//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be updated
//   - NOVEL_NOT_FOUND_ERROR if the bookmarked novel could not be updated
func (n *NovelRepository) UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error) {
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("user_id = ? AND novel_id = ?", novel.UserID, novel.NovelID).
		Update("status", novel.Status).
		Update("score", novel.Score).
		Update("current_chapter", novel.CurrentChapter).
		Error; err != nil {

		if err.Error() == "record not found" {
			return novel, types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return novel, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to update bookmarked novel", err)
	}
	return novel, nil
}

// CreateBookmarkedNovel creates a new bookmarked novel in the database.
//
// Parameters:
//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - *models.BookmarkedNovel (pointer to BookmarkedNovel struct)
//   - CONFLICT_ERROR if the bookmarked novel already exists
//   - INTERNAL_SERVER_ERROR if the bookmarked novel could not be created
func (n *NovelRepository) CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error) {
	if IsBookmarkedNovelCreated := n.isBookmarkedNovelCreated(bookmarkedNovel); IsBookmarkedNovelCreated {
		return nil, types.WrapError(types.CONFLICT_ERROR, "Bookmarked novel already exists", nil)
	}

	// Save the bookmarked novel
	if err := n.db.Create(&bookmarkedNovel).Error; err != nil {
		log.Println(err)
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to create bookmarked novel", err)
	}

	return &bookmarkedNovel, nil
}

// IsBookmarkedNovelCreated checks if a bookmarked novel with the given novel ID and user ID already exists in the database.
//
// Parameters:
//   - bookmarkedNovel models.BookmarkedNovel (BookmarkedNovel struct)
//
// Returns:
//   - bool (true if the bookmarked novel already exists, false otherwise)
func (n *NovelRepository) isBookmarkedNovelCreated(bookmarkedNovel models.BookmarkedNovel) bool {
	var existingBookmarkedNovel models.BookmarkedNovel
	if err := n.db.Where("novel_id = ? AND user_id = ?", bookmarkedNovel.NovelID, bookmarkedNovel.UserID).First(&existingBookmarkedNovel).Error; err != nil {
		return false
	}
	return existingBookmarkedNovel.ID != 0
}

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
func (n *NovelRepository) DeleteBookmarkedNovel(userID uint, novelID uint) error {
	err := n.db.Model(&models.BookmarkedNovel{}).
		Where("user_id = ? AND novel_id = ?", userID, novelID).
		Delete(&models.BookmarkedNovel{}).Error

	if err != nil {

		if err.Error() == "record not found" {
			return types.WrapError(types.NOVEL_NOT_FOUND_ERROR, "Novel not found", nil)
		}

		return types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to delete bookmarked novel", err)
	}
	return nil
}
