package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/types"
	"backend/internal/types/errors"
	"backend/internal/utils"
	"backend/internal/validators"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type NovelService struct {
	repo           interfaces.NovelRepositoryInterface
	scriptExecutor utils.ScriptExecutor
}

func NewNovelService(repo interfaces.NovelRepositoryInterface, scriptExecutor utils.ScriptExecutor) *NovelService {
	return &NovelService{repo: repo, scriptExecutor: scriptExecutor}
}

func (s *NovelService) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateAuthor(authorName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByAuthorName(authorName, page, limit)
}

func (s *NovelService) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateGenre(genreName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByGenreName(genreName, page, limit)
}

func (s *NovelService) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	err := validators.ValidateTag(tagName)
	if err != nil {
		return nil, 0, err
	}

	return s.repo.GetNovelsByTagName(tagName, page, limit)
}

func (s *NovelService) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	return s.repo.GetNovels(page, limit)
}

func (s *NovelService) GetNovelByUpdatesID(novelUpdatesID string) (*models.Novel, error) {
	novelUpdatesID, err := utils.NewNovelUpdatesIDParser().Parse(novelUpdatesID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetNovelByUpdatesID(novelUpdatesID)
}

func (s *NovelService) GetNovelByID(id uint) (*models.Novel, error) {
	return s.repo.GetNovelByID(id)
}

func (s *NovelService) CreateNovel(novelUpdatesID string) (*models.Novel, error) {
	novelUpdatesID, err := utils.NewNovelUpdatesIDParser().Parse(novelUpdatesID)
	if err != nil {
		return nil, err
	}

	// Execute the Python script
	output, err := s.scriptExecutor.ExecuteScript(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-novel", novelUpdatesID)
	if err != nil {
		return nil, types.WrapError(errors.SCRIPT_ERROR, "Failed to execute Python script: "+err.Error(), http.StatusServiceUnavailable, err)
	}

	var scriptError utils.ScriptError

	// Check if the script returned a specific error
	if json.Unmarshal(output, &scriptError) == nil {
		if scriptError.Status == 404 {
			return nil, errors.ErrNovelNotFound
		}
		return nil, types.WrapError(errors.SCRIPT_ERROR, scriptError.Error, http.StatusServiceUnavailable, nil)
	}

	// Ensure the output is unmarshaled into a valid JSON object
	var result models.ImportedNovel
	err = json.Unmarshal(output, &result)
	if err != nil {
		return nil, types.WrapError(errors.IMPORTING_NOVEL, "An error occurred while importing the novel: "+err.Error(), http.StatusInternalServerError, err)
	}

	year := strings.ReplaceAll(result.Year, "\n", "")
	status := strings.ReplaceAll(result.Status, "\n", "")
	language := strings.ReplaceAll(result.Language.Name, "\n", "")
	latestChapter, err := utils.ParseInt(result.LatestChapter)
	if err != nil {
		return nil, errors.ErrInvalidLatestChapter
	}

	novel := models.Novel{
		Title:            result.Title,
		Synopsis:         result.Synopsis,
		CoverUrl:         result.CoverUrl,
		Language:         language,
		Status:           status,
		NovelUpdatesUrl:  fmt.Sprintf("https://www.lightnovelworld.co/novel/%s", novelUpdatesID),
		NovelUpdatesID:   novelUpdatesID,
		Tags:             result.Tags,
		Authors:          result.Authors,
		Genres:           result.Genres,
		Year:             year,
		ReleaseFrequency: result.ReleaseFrequency,
		LatestChapter:    latestChapter,
	}

	// Save the novel to the database
	return s.repo.CreateNovel(novel)
}
