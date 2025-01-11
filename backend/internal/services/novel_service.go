package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
)

type NovelService struct {
	repo interfaces.NovelRepositoryInterface
}

func NewNovelService(repo interfaces.NovelRepositoryInterface) *NovelService {
	return &NovelService{repo: repo}
}

func (s *NovelService) GetNovelsByAuthorID(authorID uint, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByAuthorID(authorID, page, limit)
}

func (s *NovelService) GetNovelsByGenreID(genreID uint, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByGenreID(genreID, page, limit)
}

func (s *NovelService) GetNovelsByTagID(tagID uint, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByTagID(tagID, page, limit)
}

func (s *NovelService) GetNovelByID(id uint) (*models.Novel, error) {
	return s.repo.GetNovelByID(id)
}

func (s *NovelService) GetChapterByID(id uint) (*models.Chapter, error) {
	return s.repo.GetChapterByID(id)
}

func (s *NovelService) CreateNovel(novel models.Novel) (*models.Novel, error) {
	return s.repo.CreateNovel(novel)
}

func (s *NovelService) CreateChapters(chapters []models.Chapter) (int, error) {
	return s.repo.CreateChapters(chapters)
}

func (s *NovelService) ConvertToNovel(imported models.ImportedNovel) (*models.Novel, error) {
	return s.repo.ConvertToNovel(imported)
}

func (s *NovelService) GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error) {
	return s.repo.GetChaptersByNovelID(novelID, page, limit)
}

func (s *NovelService) CreateChapter(chapter models.Chapter) (*models.Chapter, error) {
	return s.repo.CreateChapter(chapter)
}

func (s *NovelService) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error) {
	return s.repo.GetBookmarkedNovelsByUserID(userID, page, limit)
}

func (s *NovelService) GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID uint) (models.BookmarkedNovel, error) {
	return s.repo.GetBookmarkedNovelByUserIDAndNovelID(userID, novelID)
}

func (s *NovelService) UpdateBookmarkedNovel(novel models.BookmarkedNovel) (models.BookmarkedNovel, error) {
	return s.repo.UpdateBookmarkedNovel(novel)
}

func (s *NovelService) CreateBookmarkedNovel(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error) {
	return s.repo.CreateBookmarkedNovel(bookmarkedNovel)
}

func (s *NovelService) UnbookmarkNovel(userID uint, novelID uint) error {
	return s.repo.DeleteBookmarkedNovel(userID, novelID)
}