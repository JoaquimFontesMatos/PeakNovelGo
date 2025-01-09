package interfaces

import (
	"backend/internal/models"
)

type NovelRepositoryInterface interface {
	CreateAuthor(author models.Author) (*models.Author, error)
	IsAuthorCreated(author models.Author) bool
	GetAuthorByName(name string) (*models.Author, error)
	CreateGenre(genre models.Genre) (*models.Genre, error)
	GetGenreByName(name string) (*models.Genre, error)
	IsGenreCreated(genre models.Genre) bool
	CreateTag(tag models.Tag) (*models.Tag, error)
	GetTagByName(name string) (*models.Tag, error)
	IsTagCreated(tag models.Tag) bool
	CreateNovel(novel models.Novel) (*models.Novel, error)
	IsNovelCreated(novel models.Novel) bool
	GetExistingChapterURLs(urls []string) (map[string]bool, error)
	CreateChapters(chapters []models.Chapter) (int, error)
	ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error)
	GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error)
	GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelByID(id uint) (*models.Novel, error)
	GetChapterByID(id uint) (*models.Chapter, error)
}