package repositories

import (
	"backend/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type NovelRepository struct {
	db *gorm.DB
}

func NewNovelRepository(db *gorm.DB) *NovelRepository {
	return &NovelRepository{db: db}
}

func (n *NovelRepository) CreateAuthor(author models.Author) (*models.Author, error) {
	if IsAuthorCreated := n.IsAuthorCreated(author); IsAuthorCreated {
		return nil, errors.New("author already exists")
	}

	// Save the author
	if err := n.db.Create(&author).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create author")
	}

	return &author, nil
}

func (n *NovelRepository) IsAuthorCreated(author models.Author) bool {
	var existingAuthor models.Author
	if err := n.db.Where("name = ?", author.Name).First(&existingAuthor).Error; err != nil {
		return false
	}
	return existingAuthor.ID != 0
}

func (n *NovelRepository) GetAuthorByName(name string) (*models.Author, error) {
	var author models.Author
	if err := n.db.Where("name = ?", name).First(&author).Error; err != nil {
		return nil, err
	}
	return &author, nil
}

func (n *NovelRepository) CreateGenre(genre models.Genre) (*models.Genre, error) {
	if IsGenreCreated := n.IsGenreCreated(genre); IsGenreCreated {
		return nil, errors.New("genre already exists")
	}

	// Save the genre
	if err := n.db.Create(&genre).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create genre")
	}

	return &genre, nil
}

func (n *NovelRepository) GetGenreByName(name string) (*models.Genre, error) {
	var genre models.Genre
	if err := n.db.Where("name = ?", name).First(&genre).Error; err != nil {
		return nil, err
	}
	return &genre, nil
}

func (n *NovelRepository) IsGenreCreated(genre models.Genre) bool {
	var existingGenre models.Genre
	if err := n.db.Where("name = ?", genre.Name).First(&existingGenre).Error; err != nil {
		return false
	}
	return existingGenre.ID != 0
}

func (n *NovelRepository) CreateTag(tag models.Tag) (*models.Tag, error) {
	if IsTagCreated := n.IsTagCreated(tag); IsTagCreated {
		return nil, errors.New("tag already exists")
	}

	// Save the tag
	if err := n.db.Create(&tag).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create tag")
	}

	return &tag, nil
}

func (n *NovelRepository) GetTagByName(name string) (*models.Tag, error) {
	var tag models.Tag
	if err := n.db.Where("name = ?", name).First(&tag).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (n *NovelRepository) IsTagCreated(tag models.Tag) bool {
	var existingTag models.Tag
	if err := n.db.Where("name = ?", tag.Name).First(&existingTag).Error; err != nil {
		return false
	}
	return existingTag.ID != 0
}

func (n *NovelRepository) CreateNovel(novel models.Novel) (*models.Novel, error) {
	if IsNovelCreated := n.IsNovelCreated(novel); IsNovelCreated {
		return nil, errors.New("novel already exists")

	}

	// Save the novel with relationships
	if err := n.db.Create(&novel).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create novel")
	}

	return &novel, nil
}

func (n *NovelRepository) IsNovelCreated(novel models.Novel) bool {
	var existingNovel models.Novel
	if err := n.db.Where("url = ?", novel.Url).First(&existingNovel).Error; err != nil {
		return false
	}
	return existingNovel.ID != 0
}

func (n *NovelRepository) getExistingChapterURLs(chapterNums []uint, novelID *uint) (map[string]bool, error) {
	var chapters []models.Chapter
	if err := n.db.Select("chapter_no").Where("chapter_no IN ? AND novel_id = ?", chapterNums, novelID).Find(&chapters).Error; err != nil {
		return nil, err
	}

	// Create a map for quick lookups
	existingURLs := make(map[string]bool, len(chapters))
	for _, chapter := range chapters {
		existingURLs[chapter.ChapterUrl] = true
	}
	return existingURLs, nil
}

func (n *NovelRepository) CreateChapters(chapters []models.Chapter) (int, error) {
	length := 0
	n.db.Logger = n.db.Logger.LogMode(logger.Silent) // Suppresses logs

	return length, n.db.Transaction(func(tx *gorm.DB) error {
		// Filter existing chapters
		chaptersNums := make([]uint, len(chapters))
		for i, chapter := range chapters {
			chaptersNums[i] = chapter.ChapterNo
		}

		if len(chapters) == 0 {
			return errors.New("no chapters to save")
		}
		novelID := chapters[0].NovelID

		existingURLs, err := n.getExistingChapterURLs(chaptersNums, novelID)
		if err != nil {
			log.Println(err)
			return errors.New("failed to fetch existing chapters")
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
			return errors.New("no new chapters to save")
		}

		if err := tx.Create(&newChapters).Error; err != nil {
			return err
		}
		length = len(newChapters)

		log.Printf("%d chapters added to the database", length)

		return nil
	})
}

func (n *NovelRepository) getTagsByName(names []string) (map[string]models.Tag, error) {
	var existingTags []models.Tag
	if err := n.db.Where("name IN ?", names).Find(&existingTags).Error; err != nil {
		return nil, err
	}

	// Create a map for quick lookups
	tagMap := make(map[string]models.Tag, len(existingTags))
	for _, tag := range existingTags {
		tagMap[tag.Name] = tag
	}
	return tagMap, nil
}

func (n *NovelRepository) getAuthorsByName(names []string) (map[string]models.Author, error) {
	var existingAuthors []models.Author
	if err := n.db.Where("name IN ?", names).Find(&existingAuthors).Error; err != nil {
		return nil, err
	}

	// Create a map for quick lookups
	authorMap := make(map[string]models.Author, len(existingAuthors))
	for _, author := range existingAuthors {
		authorMap[author.Name] = author
	}
	return authorMap, nil
}

func (n *NovelRepository) getGenresByName(names []string) (map[string]models.Genre, error) {
	var existingGenres []models.Genre
	if err := n.db.Where("name IN ?", names).Find(&existingGenres).Error; err != nil {
		return nil, err
	}

	// Create a map for quick lookups
	genreMap := make(map[string]models.Genre, len(existingGenres))
	for _, genre := range existingGenres {
		genreMap[genre.Name] = genre
	}
	return genreMap, nil
}

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

func (n *NovelRepository) GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error) {
	var chapters []models.Chapter
	var total int64

	// Count total chapters for the novel
	if err := n.db.Model(&models.Chapter{}).Where("novel_id = ?", novelID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Where("novel_id = ?", novelID).
		Order("chapter_no ASC").
		Limit(limit).
		Offset(offset).
		Find(&chapters).Error; err != nil {
		return nil, 0, err
	}

	return chapters, total, nil
}

func (n *NovelRepository) GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the genre
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_genres ON novel_genres.novel_id = novels.id").
		Where("novel_genres.genre_id = ?", genreID).
		Count(&total).Error; err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	return novels, total, nil
}

func (n *NovelRepository) GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the tag
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_tags ON novel_tags.novel_id = novels.id").
		Where("novel_tags.tag_id = ?", tagID).
		Count(&total).Error; err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}

	return novels, total, nil
}

func (n *NovelRepository) GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error) {
	var novels []models.Novel
	var total int64

	// Count total novels for the author
	if err := n.db.Model(&models.Novel{}).
		Joins("JOIN novel_authors ON novel_authors.novel_id = novels.id").
		Where("novel_authors.author_id = ?", authorID).
		Count(&total).Error; err != nil {
		return nil, 0, err
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
		return nil, 0, err
	}
	return novels, total, nil
}

func (n *NovelRepository) GetNovelByID(id uint) (*models.Novel, error) {
	var novel models.Novel
	if err := n.db.Where("id = ?", id).
		Preload("Authors").
		Preload("Genres").
		Preload("Tags").
		First(&novel).
		Error; err != nil {
		return nil, err
	}
	return &novel, nil
}

func (n *NovelRepository) GetChapterByID(id uint) (*models.Chapter, error) {
	var chapter models.Chapter
	if err := n.db.Where("id = ?", id).First(&chapter).Error; err != nil {
		return nil, err
	}
	return &chapter, nil
}

func (n *NovelRepository) IsChapterCreated(chapter models.Chapter) bool {
	var existingChapter models.Chapter
	if err := n.db.Where("chapter_no = ? AND novel_id = ?", chapter.ChapterNo, chapter.NovelID).First(&existingChapter).Error; err != nil {
		return false
	}
	return existingChapter.ID != 0
}

func (n *NovelRepository) CreateChapter(chapter models.Chapter) (*models.Chapter, error) {
	if IsChapterCreated := n.IsChapterCreated(chapter); IsChapterCreated {
		return nil, errors.New("chapter already exists")
	}

	// Save the chapter
	if err := n.db.Create(&chapter).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create chapter")
	}

	return &chapter, nil
}

func (n *NovelRepository) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error) {
	var novels []models.BookmarkedNovel
	var total int64

	// Count total novels for the user
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ?", userID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ?", userID).
		Limit(limit).
		Offset(offset).
		Find(&novels).Error; err != nil {
		return nil, 0, err
	}

	return novels, total, nil
}

func (n *NovelRepository) GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error) {
	var novel models.BookmarkedNovel
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("bookmarked_novels.user_id = ? AND bookmarked_novels.novel_id = ?", userID, novelID).
		First(&novel).Error; err != nil {
		return novel, err
	}
	return novel, nil
}

func (n *NovelRepository) UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error) {
	if err := n.db.Model(&models.BookmarkedNovel{}).
		Where("user_id = ? AND novel_id = ?", novel.UserID, novel.NovelID).
		Update("status", novel.Status).
		Update("score", novel.Score).
		Update("current_chapter", novel.CurrentChapter).
		Error; err != nil {
		return novel, err
	}
	return novel, nil
}

func (n *NovelRepository) CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error) {
	if IsBookmarkedNovelCreated := n.IsBookmarkedNovelCreated(bookmarkedNovel); IsBookmarkedNovelCreated {
		return nil, errors.New("bookmarked novel already exists")
	}

	// Save the bookmarked novel
	if err := n.db.Create(&bookmarkedNovel).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create bookmarked novel")
	}

	return &bookmarkedNovel, nil
}

func (n *NovelRepository) IsBookmarkedNovelCreated(bookmarkedNovel models.BookmarkedNovel) bool {
	var existingBookmarkedNovel models.BookmarkedNovel
	if err := n.db.Where("novel_id = ? AND user_id = ?", bookmarkedNovel.NovelID, bookmarkedNovel.UserID).First(&existingBookmarkedNovel).Error; err != nil {
		return false
	}
	return existingBookmarkedNovel.ID != 0
}
