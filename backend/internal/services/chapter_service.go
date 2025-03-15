package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/types/errors"
	"backend/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

type ChapterService struct {
	repo interfaces.ChapterRepositoryInterface
}

func NewChapterService(repo interfaces.ChapterRepositoryInterface) *ChapterService {
	return &ChapterService{repo: repo}
}

func (s *ChapterService) IsChapterCreated(chapterNo uint, novelID uint) bool {
	return s.repo.IsChapterCreated(chapterNo, novelID)
}

func (s *ChapterService) CreateChapter(novelID uint, result models.ImportedChapterMetadata) error {
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

func (s *ChapterService) ImportChapter(novelUpdatesID string, chapterNo int) (models.ImportedChapterMetadata, error) {
	chapterNoStr := strconv.Itoa(chapterNo)
	cmd := exec.Command(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-chapter", novelUpdatesID, chapterNoStr)

	output, err := cmd.Output()
	if err != nil {
		return models.ImportedChapterMetadata{}, types.WrapError(errors.SCRIPT_ERROR, "Failed to execute Python script: "+err.Error(), http.StatusServiceUnavailable, err)
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
		return models.ImportedChapterMetadata{}, types.WrapError(errors.IMPORTING_CHAPTER, "Failed to parse Python script output as JSON", http.StatusInternalServerError, err)
	}

	if result.Status == 404 {
		return models.ImportedChapterMetadata{}, errors.ErrNoChapters
	}

	if result.Status == 204 {
		return models.ImportedChapterMetadata{}, errors.ErrChapterNotFound
	}

	return models.ImportedChapterMetadata{
		ID:         uint(chapterNo),
		Title:      result.Title,
		Body:       utils.StripHTML(result.Body),
		ChapterUrl: fmt.Sprintf("https://www.lightnovelworld.co/novel/%s/chapter-%d", novelUpdatesID, chapterNo),
	}, nil
}

func (s *ChapterService) GetChapterByNovelUpdatesIDAndChapterNo(novelTitle string, chapterNo uint) (*models.Chapter, error) {
	return s.repo.GetChapterByNovelUpdatesIDAndChapterNo(novelTitle, chapterNo)
}

func (s *ChapterService) GetChaptersByNovelUpdatesID(novelTitle string, page, limit int) ([]models.Chapter, int64, error) {
	return s.repo.GetChaptersByNovelUpdatesID(novelTitle, page, limit)
}
