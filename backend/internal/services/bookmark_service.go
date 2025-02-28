package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
)

type BookmarkService struct {
	repo interfaces.BookmarkRepositoryInterface
}

func NewBookmarkService(repo interfaces.BookmarkRepositoryInterface) *BookmarkService {
	return &BookmarkService{repo: repo}
}

func (s *BookmarkService) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetBookmarkedNovelsByUserID(userID, page, limit)
}

func (s *BookmarkService) GetBookmarkByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error) {
	return s.repo.GetBookmarkByUserIDAndNovelID(userID, novelID)
}

func (s *BookmarkService) UpdateBookmark(novel models.BookmarkedNovel) (models.BookmarkedNovel, error) {
	return s.repo.UpdateBookmark(novel)
}

func (s *BookmarkService) CreateBookmark(bookmarkedNovel models.BookmarkedNovel) (*models.BookmarkedNovel, error) {
	return s.repo.CreateBookmark(bookmarkedNovel)
}

func (s *BookmarkService) UnbookmarkNovel(userID uint, novelID uint) error {
	return s.repo.DeleteBookmark(userID, novelID)
}
