package controllers

import (
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/utils"
	"backend/internal/validators"
	"strings"

	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NovelController struct {
	novelService interfaces.NovelServiceInterface
}

func NewNovelController(novelService interfaces.NovelServiceInterface) *NovelController {
	return &NovelController{novelService: novelService}
}

// HandleGetNovels handles POST /novel
func (n *NovelController) HandleImportNovel(ctx *gin.Context) {
	var metadata models.ImportedNovel

	if err := validators.ValidateBody(ctx, &metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	year := strings.ReplaceAll(metadata.Year, "\n", "")
	status := strings.ReplaceAll(metadata.Status, "\n", "")
	language := strings.ReplaceAll(metadata.Language.Name, "\n", "")
	lowerCaseTitle := strings.ToLower(metadata.Title)
	novelUpdatesID := strings.ReplaceAll(lowerCaseTitle, " ", "-")

	novel := models.Novel{
		Title:            metadata.Title,
		Synopsis:         metadata.Synopsis,
		CoverUrl:         metadata.CoverUrl,
		Language:         language,
		Status:           status,
		NovelUpdatesUrl:  fmt.Sprintf("https://www.novelupdates.com/series/%s", novelUpdatesID),
		NovelUpdatesID:   novelUpdatesID,
		Tags:             metadata.Tags,
		Authors:          metadata.Authors,
		Genres:           metadata.Genres,
		Year:             year,
		ReleaseFrequency: metadata.ReleaseFrequency,
	}

	// Save the novel to the database
	createdNovel, err := n.novelService.CreateNovel(novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save novel to database"})
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(200, createdNovel)
}

// HandleGetNovels handles POST /novel
func (n *NovelController) HandleImportNovelByNovelUpdatesID(ctx *gin.Context) {
	novelUpdatesID := ctx.Param("novel_updates_id")

	if novelUpdatesID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No novel updates ID provided"})
		return
	}

	// Specify the Python script and its module
	cmd := exec.Command(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-novel", novelUpdatesID)

	// Capture the output of the Python script
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error: ", err)
		ctx.JSON(500, gin.H{"error": "Failed to execute Python script"})
		return
	}

	// Ensure the output is unmarshaled into a valid JSON object
	var result models.ImportedNovel
	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		ctx.JSON(500, gin.H{"error": "Failed to parse Python script output as JSON"})
		return
	}

	year := strings.ReplaceAll(result.Year, "\n", "")
	status := strings.ReplaceAll(result.Status, "\n", "")
	language := strings.ReplaceAll(result.Language.Name, "\n", "")
	latestChapter, err := utils.ParseInt(result.LatestChapter)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Invalid Latest Chapter"})
		return
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
	createdNovel, err := n.novelService.CreateNovel(novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save novel to database"})
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(200, createdNovel)

}

func (n *NovelController) GetNovelsByAuthorName(ctx *gin.Context) {
	authorName := ctx.Param("author_name")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByAuthorName(authorName, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Build response with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"data":       novels,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (n *NovelController) GetNovelsByGenreName(ctx *gin.Context) {
	genreName := ctx.Param("genre_name")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByGenreName(genreName, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Build response with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"data":       novels,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (n *NovelController) GetNovelsByTagName(ctx *gin.Context) {
	tagName := ctx.Param("tag_name")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByTagName(tagName, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Build response with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"data":       novels,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (n *NovelController) GetNovels(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovels(page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Build response with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"data":       novels,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (n *NovelController) GetNovelByID(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novel, err := n.novelService.GetNovelByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

func (n *NovelController) GetNovelByUpdatesID(ctx *gin.Context) {
	title := ctx.Param("title")

	novel, err := n.novelService.GetNovelByUpdatesID(title)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}


