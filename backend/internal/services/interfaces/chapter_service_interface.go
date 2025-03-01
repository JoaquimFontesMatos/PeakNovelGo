package interfaces

import (
	"backend/internal/models"
)

type ChapterServiceInterface interface {
	IsChapterCreated(chapterNo uint, novelID uint) bool
	CreateChapter(novelID uint, result models.ImportedChapterMetadata) error
	ImportChapter(novelUpdatesID string, chapterNo int) (models.ImportedChapterMetadata, error)
	GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error)
	GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error)
}