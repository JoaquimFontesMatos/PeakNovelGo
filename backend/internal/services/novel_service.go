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

func (s *NovelService) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByAuthorName(authorName, page, limit)
}

func (s *NovelService) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByGenreName(genreName, page, limit)
}

func (s *NovelService) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovelsByTagName(tagName, page, limit)
}

func (s *NovelService) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovels(page, limit)
}

func (s *NovelService) GetNovelByUpdatesID(title string) (*models.Novel, error) {
	return s.repo.GetNovelByUpdatesID(title)
}

func (s *NovelService) GetNovelByID(id uint) (*models.Novel, error) {
	return s.repo.GetNovelByID(id)
}

func (s *NovelService) CreateNovel(novel models.Novel) (*models.Novel, error) {
	return s.repo.CreateNovel(novel)
}
