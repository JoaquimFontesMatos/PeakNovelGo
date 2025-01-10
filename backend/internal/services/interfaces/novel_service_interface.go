package interfaces

import (
	"backend/internal/models"
)

type NovelServiceInterface interface {
	GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelByID(id uint) (*models.Novel, error)
	GetChapterByID(id uint) (*models.Chapter, error)
	CreateNovel(novel models.Novel) (*models.Novel, error)
	CreateChapters(chapters []models.Chapter) (int, error)
	ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error)
	GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error)
}
