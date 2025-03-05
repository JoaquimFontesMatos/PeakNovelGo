package interfaces

import (
	"backend/internal/models"
)

type NovelServiceInterface interface {
	GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error)
	GetNovels(page, limit int) ([]models.Novel, int64, error)
	GetNovelByID(id uint) (*models.Novel, error)
	GetNovelByUpdatesID(title string) (*models.Novel, error)
	CreateNovel(novelUpdatesID string) (*models.Novel, error)
}
