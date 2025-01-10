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
	CreateChapters(chapters []models.Chapter) (int, error)
	ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error)
	GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error)
	GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error)
	GetNovelByID(id uint) (*models.Novel, error)
	GetChapterByID(id uint) (*models.Chapter, error)
	IsChapterCreated(chapter models.Chapter) bool
	CreateChapter(chapter models.Chapter) (*models.Chapter, error)
	GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error)
	GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error)
	UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error)
	CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error)
	IsBookmarkedNovelCreated(bookmarkedNovel models.BookmarkedNovel) bool
}