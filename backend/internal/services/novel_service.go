package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
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

func (s *NovelService) GetChapterByID(id uint) (*models.Chapter, error) {
	return s.repo.GetChapterByID(id)
}

func (s *NovelService) GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error) {
	return s.repo.GetChapterByNovelUpdatesIDAndChapterNo(novelTitle, chapterNo)
}

func (s *NovelService) GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error) {
	return s.repo.GetChaptersByNovelUpdatesID(novelTitle, page, limit)
}

func (s *NovelService) CreateNovel(novel models.Novel) (*models.Novel, error) {
	return s.repo.CreateNovel(novel)
}

func (s *NovelService) CreateChapters(chapters []models.Chapter) (int, error) {
	return s.repo.CreateChapters(chapters)
}

func (s *NovelService) GetChaptersByNovelID(novelID uint, page, limit int) ([]models.Chapter, int64, error) {
	return s.repo.GetChaptersByNovelID(novelID, page, limit)
}

func (s *NovelService) IsChapterCreated(chapterNo uint, novelID uint) bool {
	return s.repo.IsChapterCreated(chapterNo, novelID)
}

func (s *NovelService) CreateChapter(novelID uint, result models.ImportedChapterMetadata) error {
	importedChapter := models.ImportedChapter{
		NovelID:    &novelID,
		ID:         result.ID,
		Title:      result.Title,
		ChapterUrl: result.ChapterUrl,
		Body:       result.Body,
	}

	chapter := importedChapter.ToChapter()
	_, err := s.repo.CreateChapter(*chapter)
	if err != nil {
		return err
	}

	return nil
}

func (s *NovelService) ImportChapter(novelUpdatesID string, chapterNo int) (models.ImportedChapterMetadata, error) {
	chapterNoStr := strconv.Itoa(chapterNo)
	cmd := exec.Command(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-chapter", novelUpdatesID, chapterNoStr)

	output, err := cmd.Output()
	if err != nil {
		return models.ImportedChapterMetadata{}, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to execute Python script", err)
	}

	var result struct {
		Title     string `json:"title"`
		Body      string `json:"body"`
		ChapterNo string `json:"chapter_no"`
		Status    int    `json:"status"`
		Error     string `json:"error,omitempty"`
	}

	err = json.Unmarshal(output, &result)
	if err != nil {
		return models.ImportedChapterMetadata{}, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to parse Python script output as JSON", err)
	}

	if result.Status == 404 {
		return models.ImportedChapterMetadata{}, types.WrapError(types.NO_CHAPTERS_ERROR, "No more chapters available", nil)
	}

	if result.Status == 204 {
		return models.ImportedChapterMetadata{},types.WrapError(types.CHAPTER_NOT_FOUND_ERROR, "Skipping empty chapter", nil)
	}

	return models.ImportedChapterMetadata{
		ID:         uint(chapterNo),
		Title:      result.Title,
		Body:       utils.StripHTML(result.Body),
		ChapterUrl: fmt.Sprintf("https://www.lightnovelworld.co/novel/%s/chapter-%d", novelUpdatesID, chapterNo),
	}, nil
}

func (s *NovelService) GetBookmarkedNovelsByUserID(userID uint, page, limit int) ([]models.BookmarkedNovel, int64, error) {
	return s.repo.GetBookmarkedNovelsByUserID(userID, page, limit)
}

func (s *NovelService) GetBookmarkedNovelByUserIDAndNovelID(userID uint, novelID string) (models.BookmarkedNovel, error) {
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
