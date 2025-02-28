package controllers

import (
	"backend/internal/dtos"
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

type ChapterController struct {
	chapterService interfaces.ChapterServiceInterface
	novelService   interfaces.NovelServiceInterface
}

func NewChapterController(chapterService interfaces.ChapterServiceInterface, novelService interfaces.NovelServiceInterface) *ChapterController {
	return &ChapterController{chapterService: chapterService, novelService: novelService}
}

// HandleImportChapters handles streaming response for importing chapters
func (c *ChapterController) HandleImportChapters(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")

	novel, err := c.novelService.GetNovelByUpdatesID(idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush() // Ensure headers are sent immediately

	err = c.processChaptersWithStreaming(ctx, idParam, novel.ID)
	if err != nil {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}
}

func (c *ChapterController) processChaptersWithStreaming(ctx *gin.Context, novelUpdatesID string, novelID uint) error {
	const workerCount = 10
	chapterCount := 0
	chapterQueue := make(chan int, workerCount)
	results := make(chan models.ImportedChapterMetadata)
	errorsChan := make(chan dtos.ChapterStatus)
	skippedChan := make(chan dtos.ChapterStatus)
	done := make(chan struct{})

	var wg sync.WaitGroup

	// Fetch the latest chapter number before scraping
	seriesInfo, err := c.novelService.GetNovelByUpdatesID(novelUpdatesID)
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
				isChapterCreated := c.chapterService.IsChapterCreated(uint(chapterNo), novelID)
				if isChapterCreated {
					skippedChan <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "skipped"}
					continue
				}

				result, err := c.chapterService.ImportChapter(novelUpdatesID, chapterNo)
				if err != nil {
					errorsChan <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "error"}
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

		case result, ok := <-results:
			if !ok {
				log.Printf("All chapters processed")
				fmt.Fprintf(ctx.Writer, "event: complete\ndata: All chapters processed\n\n")
				ctx.Writer.Flush()
				return nil
			}

			err := c.chapterService.CreateChapter(novelID, result)
			if err != nil {
				log.Printf("Failed to save chapter %d: %v", result.ID, err)
				return err
			}

			chapterCount++
			chapterStatuses[int(result.ID)] = "downloaded"

		case skipped, ok := <-skippedChan:
			if ok {
				chapterStatuses[skipped.ChapterNo] = "skipped"
			}
		}

		fmt.Fprintf(ctx.Writer, "event: status\ndata: %s\n\n", utils.GetStatusJSON(chapterStatuses))
		ctx.Writer.Flush()
	}
}

func (c *ChapterController) GetChapterByNovelUpdatesIDAndChapterNo(ctx *gin.Context) {
	novelTitle := ctx.Param("novel_title")
	chapterNo := ctx.Param("chapter_no")

	chapterNoUint, err := utils.ParseID(chapterNo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	chapter, err := c.chapterService.GetChapterByNovelUpdatesIDAndChapterNo(novelTitle, chapterNoUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, chapter)
}

func (c *ChapterController) GetChaptersByNovelUpdatesID(ctx *gin.Context) {
	novelTitle := ctx.Param("novel_title")

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	chapters, total, err := c.chapterService.GetChaptersByNovelUpdatesID(novelTitle, page, limit)
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
