package controllers

import (
	"archive/zip"
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"backend/internal/validators"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type NovelController struct {
	novelService interfaces.NovelServiceInterface
}

func NewNovelController(novelService interfaces.NovelServiceInterface) *NovelController {
	return &NovelController{novelService: novelService}
}

type Session struct {
	UserInput        string `json:"user_input,omitempty"`
	OutputPath       string `json:"output_path,omitempty"`
	Completed        bool   `json:"completed,omitempty"`
	DownloadChapters []int  `json:"download_chapters,omitempty"`
}

type Metadata struct {
	Novel   models.ImportedNovel `json:"novel,omitempty"`
	Session Session              `json:"session,omitempty"`
}

// HandleGetNovels handles POST /novel
func (n *NovelController) HandleImportNovel(ctx *gin.Context) {
	if ctx.Request.Body == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No body provided"})
		return
	}

	metadata := Metadata{}

	if err := validators.ValidateBody(ctx, &metadata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	metadata.Novel.Synopsis = utils.StripHTML(metadata.Novel.Synopsis)

	// Convert the imported novel to a Novel struct
	novel, err := n.novelService.ConvertToNovel(metadata.Novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert imported novel to Novel struct"})
		return
	}

	// Save the novel to the database
	createdNovel, err := n.novelService.CreateNovel(*novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save novel to database"})
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(200, createdNovel)
}

// HandleUploadNovelZip handles POST /novel/upload
func (n *NovelController) HandleImportChaptersZip(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Save the uploaded file to a temporary location
	tempFile := "./temp_upload.zip"
	if err := ctx.SaveUploadedFile(file, tempFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Ensure the temp file is deleted after processing, even if an error occurs
	defer func() {
		if err := os.Remove(tempFile); err != nil {
			// Log the error if the file cannot be deleted
			log.Printf("Failed to delete temp file: %v\n", err)
		}
	}()

	// Process the ZIP file
	chapterCount, err := n.processChaptersZip(tempFile, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the extracted metadata
	string := fmt.Sprintf("Chapters extracted successfully. %d chapters added.", chapterCount)
	ctx.JSON(http.StatusOK, gin.H{"message": string})
}

func (n *NovelController) processChaptersZip(filePath string, uid uint) (int, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	chapters := make([]models.Chapter, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	errChan := make(chan error, len(reader.File))

	// Process each file concurrently
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		wg.Add(1)
		go func(file *zip.File) {
			defer wg.Done()

			f, err := file.Open()

			if err != nil {
				errChan <- fmt.Errorf("failed to open file in zip: %v", err)
				return
			}

			defer f.Close()

			var chapterData models.ImportedChapter
			if err := json.NewDecoder(f).Decode(&chapterData); err != nil {
				errChan <- fmt.Errorf("failed to decode JSON file: %v", err)
				return
			}

			// Clean up the chapter body
			chapterData.Body = utils.StripHTML(chapterData.Body)
			chapterData.NovelID = &uid

			chapter := chapterData.ToChapter()

			mu.Lock()
			chapters = append(chapters, *chapter)
			mu.Unlock()
		}(file)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return 0, err
		}
	}

	// Save the chapters to the database
	return n.novelService.CreateChapters(chapters)
}

func (n *NovelController) GetChaptersByNovelID(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse query parameters
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	chapters, total, err := n.novelService.GetChaptersByNovelID(id, page, limit)
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
		"data":       chapters,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func (n *NovelController) GetNovelsByAuthorID(ctx *gin.Context) {
	idParam := ctx.Param("author_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByAuthorID(id, page, limit)
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

func (n *NovelController) GetNovelsByGenreID(ctx *gin.Context) {
	idParam := ctx.Param("genre_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByGenreID(id, page, limit)
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

func (n *NovelController) GetNovelsByTagID(ctx *gin.Context) {
	idParam := ctx.Param("tag_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetNovelsByTagID(id, page, limit)
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

func (n *NovelController) GetChapterByID(ctx *gin.Context) {
	idParam := ctx.Param("chapter_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := n.novelService.GetChapterByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chapter)
}

func (n *NovelController) CreateChapter(ctx *gin.Context) {
	var chapter models.Chapter
	if err := ctx.ShouldBindJSON(&chapter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdChapter, err := n.novelService.CreateChapter(chapter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, createdChapter)
}

func (n *NovelController) GetBookmarkedNovelsByUserID(ctx *gin.Context) {
	idParam := ctx.Param("user_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	novels, total, err := n.novelService.GetBookmarkedNovelsByUserID(id, page, limit)
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

func (n *NovelController) CreateBookmarkedNovel(ctx *gin.Context) {
	var bookmarkedNovel models.BookmarkedNovel
	if err := ctx.ShouldBindJSON(&bookmarkedNovel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBookmarkedNovel, err := n.novelService.CreateBookmarkedNovel(bookmarkedNovel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save bookmarked novel to database"})
		return
	}

	log.Println("Bookmarked novel saved successfully")

	ctx.JSON(http.StatusOK, createdBookmarkedNovel)
}

func (n *NovelController) UpdateBookmarkedNovel(ctx *gin.Context) {
	var novel models.BookmarkedNovel
	if err := ctx.ShouldBindJSON(&novel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNovel, err := n.novelService.UpdateBookmarkedNovel(novel)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedNovel)
}

func (n *NovelController) GetBookmarkedNovelByUserIDAndNovelID(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")
	novelIDParam := ctx.Param("novel_id")

	userID, err := utils.ParseID(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	novelID, err := utils.ParseID(novelIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novel, err := n.novelService.GetBookmarkedNovelByUserIDAndNovelID(userID, novelID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

func (n *NovelController) UnbookmarkNovel(ctx *gin.Context) {
	userID, err := strconv.ParseUint(ctx.Param("userID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	novelID, err := strconv.ParseUint(ctx.Param("novelID"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = n.novelService.UnbookmarkNovel(uint(userID), uint(novelID))
	if err != nil {
		error := err.(*types.MyError)
		if error.Code == types.NOVEL_NOT_FOUND_ERROR {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Novel successfully unbookmarked"})
}