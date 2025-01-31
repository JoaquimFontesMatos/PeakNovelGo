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
	GetChapterByID(id uint) (*models.Chapter, error)
	GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error)
	GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error)
	CreateNovel(novel models.Novel) (*models.Novel, error)
	CreateChapters(chapters []models.Chapter) (int, error)
	GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error)
	CreateChapter(chapter models.Chapter) (*models.Chapter, error)
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error)
	GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error)
	UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error)
	CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)
	UnbookmarkNovel(userID uint, novelID uint) error
}
