package controllers

import (
	"archive/zip"
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"backend/internal/validators"
	"errors"
	"strings"

	//"backend/internal/validators"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
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
		NovelUpdatesUrl:  fmt.Sprintf("https://www.novelupdates.com/series/%s", novelUpdatesID),
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

// HandleImportChapters handles streaming response for importing chapters
func (n *NovelController) HandleImportChapters(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")

	novel, err := n.novelService.GetNovelByUpdatesID(idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush() // Ensure headers are sent immediately

	err = n.processChaptersWithStreaming(ctx, idParam, novel.ID)
	if err != nil {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}
}

type ChapterStatus struct {
	ChapterNo int    `json:"chapterNo"`
	Status    string `json:"status"`
}

func (n *NovelController) processChaptersWithStreaming(ctx *gin.Context, novelUpdatesID string, novelID uint) error {
	const workerCount = 10
	chapterCount := 0
	chapterQueue := make(chan int, workerCount)
	results := make(chan models.ImportedChapterMetadata)
	errorsChan := make(chan ChapterStatus)
	skippedChan := make(chan ChapterStatus)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// Fetch the latest chapter number before scraping
	seriesInfo, err := n.novelService.GetNovelByUpdatesID(novelUpdatesID)
	if err != nil {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return err
	}

	latestChapter := seriesInfo.LatestChapter
	if latestChapter == 0 {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: No chapters found\n\n")
		ctx.Writer.Flush()
		return fmt.Errorf("no chapters found")
	}

	chapterStatuses := make(map[int]string, latestChapter+2)
	for i := 1; i <= latestChapter; i++ {
		chapterStatuses[i] = "to download"
	}

	// Worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chapterNo := range chapterQueue {
				isChapterCreated := n.novelService.IsChapterCreated(uint(chapterNo), novelID)
				if isChapterCreated {
					skippedChan <- ChapterStatus{ChapterNo: chapterNo, Status: "skipped"}
					continue
				}

				result, err := n.importChapter(novelUpdatesID, chapterNo)
				if err != nil {
					errorsChan <- ChapterStatus{ChapterNo: chapterNo, Status: "error"}
					continue
				}

				results <- result
			}
		}()
	}

	// Populate queue up to the latest chapter
	go func() {
		defer close(chapterQueue)
		for chapterNo := 1; chapterNo <= latestChapter; chapterNo++ {
			select {
			case <-done:
				return
			case chapterQueue <- chapterNo:
				log.Printf("Added chapter %d to queue", chapterNo)
			}
		}
	}()

	// Goroutine to close result-related channels after workers finish
	go func() {
		wg.Wait()
		close(results)
		close(errorsChan)
		close(skippedChan)
	}()

	// Process results
	for {
		select {
		case err := <-errorsChan:
			chapterStatuses[err.ChapterNo] = err.Status
			log.Printf("Error importing chapter %d: %s", err.ChapterNo, err.Status)
			ctx.Writer.Flush()

		case result, ok := <-results:
			if !ok {
				log.Printf("All chapters processed")
				fmt.Fprintf(ctx.Writer, "event: complete\ndata: All chapters processed\n\n")
				ctx.Writer.Flush()
				return nil
			}

			chapterStatuses[int(result.ID)] = "downloading"
			ctx.Writer.Flush()

			err := n.saveChapter(novelID, result)
			if err != nil {
				log.Printf("Failed to save chapter %d: %v", result.ID, err)
				return err
			}

			chapterCount++
			chapterStatuses[int(result.ID)] = "downloaded"
			ctx.Writer.Flush()

		case skipped, ok := <-skippedChan:
			if ok {
				chapterStatuses[skipped.ChapterNo] = "skipped"
				ctx.Writer.Flush()
			}
		}

		// Send the current status map
		statusJSON, err := json.Marshal(chapterStatuses)
		if err != nil {
			log.Printf("Failed to marshal statuses: %v", err)
			fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
			ctx.Writer.Flush()
			return err
		}

		fmt.Fprintf(ctx.Writer, "event: status\ndata: %s\n\n", statusJSON)
		ctx.Writer.Flush()
	}
}

func (n *NovelController) importChapter(novelUpdatesID string, chapterNo int) (models.ImportedChapterMetadata, error) {
	chapterNoStr := strconv.Itoa(chapterNo)
	cmd := exec.Command(os.Getenv("PYTHON"), "-m", "novel_updates_scraper.client", "import-chapter", novelUpdatesID, chapterNoStr)

	output, err := cmd.Output()
	if err != nil {
		return models.ImportedChapterMetadata{}, fmt.Errorf("failed to execute Python script: %v", err)
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
		return models.ImportedChapterMetadata{}, fmt.Errorf("failed to parse Python script output as JSON: %v", err)
	}

	if result.Status == 404 {
		return models.ImportedChapterMetadata{}, fmt.Errorf("no more chapters available")
	}

	if result.Status == 204 {
		return models.ImportedChapterMetadata{}, fmt.Errorf("skipping empty chapter %d", chapterNo)
	}

	return models.ImportedChapterMetadata{
		ID:         uint(chapterNo),
		Title:      result.Title,
		Body:       utils.StripHTML(result.Body),
		ChapterUrl: fmt.Sprintf("https://www.lightnovelworld.co/novel/%s/chapter-%d", novelUpdatesID, chapterNo),
	}, nil
}

// saveChapter inserts chapter into database
func (n *NovelController) saveChapter(novelID uint, result models.ImportedChapterMetadata) error {
	importedChapter := models.ImportedChapter{
		NovelID:    &novelID,
		ID:         result.ID,
		Title:      result.Title,
		ChapterUrl: result.ChapterUrl,
		Body:       result.Body,
	}

	chapter := importedChapter.ToChapter()
	_, err := n.novelService.CreateChapter(*chapter)
	if err != nil {
		var userErr *types.MyError
		if errors.As(err, &userErr) {
			if userErr.Code == "CONFLICT_ERROR" {
				return nil // Ignore duplicate errors
			}
		}
		return err
	}

	return nil
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

func (n *NovelController) GetChapterByNovelUpdatesIDAndChapterNo(ctx *gin.Context) {
	novelTitle := ctx.Param("novel_title")
	chapterNo := ctx.Param("chapter_no")

	chapterNoUint, err := utils.ParseID(chapterNo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := n.novelService.GetChapterByNovelUpdatesIDAndChapterNo(novelTitle, chapterNoUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chapter)
}

func (n *NovelController) GetChaptersByNovelUpdatesID(ctx *gin.Context) {
	novelTitle := ctx.Param("novel_title")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	chapters, total, err := n.novelService.GetChaptersByNovelUpdatesID(novelTitle, page, limit)
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

	novel, err := n.novelService.GetBookmarkedNovelByUserIDAndNovelID(userID, novelIDParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

func (n *NovelController) UnbookmarkNovel(ctx *gin.Context) {

	userID, err := utils.ParseID(ctx.Param("user_id"))

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	novelID, err := utils.ParseID(ctx.Param("novel_id"))

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
