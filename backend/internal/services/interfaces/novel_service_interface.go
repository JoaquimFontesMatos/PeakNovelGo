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
	CreateChapter(chapter models.Chapter) (*models.Chapter, error)
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error)
	GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error)
	UpdateBookmarkedNovel( novel models.BookmarkedNovel) (models.BookmarkedNovel, error)
	CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)
	UnbookmarkNovel(userID uint, novelID uint) error
}
