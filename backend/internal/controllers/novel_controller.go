package controllers

import (
	"archive/zip"
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
	"backend/internal/utils"
	"backend/internal/validators"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type NovelController struct {
	novelRepository interfaces.NovelRepositoryInterface
}

func NewNovelController(novelRepository interfaces.NovelRepositoryInterface) *NovelController {
	return &NovelController{novelRepository: novelRepository}
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
	novel, err := n.novelRepository.ConvertToNovel(metadata.Novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to convert imported novel to Novel struct"})
		return
	}

	// Save the novel to the database
	createdNovel, err := n.novelRepository.CreateNovel(*novel)

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
	id, err := validators.ValidateID(idParam)
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

			var chapterData struct {
				Title      string `json:"title"`
				Body       string `json:"body"`
				ChapterUrl string `json:"url"`
			}
			if err := json.NewDecoder(f).Decode(&chapterData); err != nil {
				errChan <- fmt.Errorf("failed to decode JSON file: %v", err)
				return
			}

			f.Close()

			// Clean up the chapter body
			chapterData.Body = utils.StripHTML(chapterData.Body)

			chapter := models.Chapter{
				NovelID:    &uid,
				Title:      chapterData.Title,
				Body:       chapterData.Body,
				ChapterUrl: chapterData.ChapterUrl,
			}

			mu.Lock()
			chapters = append(chapters, chapter)
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
	return n.novelRepository.CreateChapters(chapters)
}

func (n *NovelController) GetChaptersByNovelID(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := validators.ValidateID(idParam)
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

	chapters, total, err := n.novelRepository.GetChaptersByNovelID(id, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Build response with pagination metadata
	ctx.JSON(http.StatusOK, gin.H{
		"data":       chapters,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": (total + int64(limit) - 1) / int64(limit), // Calculate total pages
	})
}