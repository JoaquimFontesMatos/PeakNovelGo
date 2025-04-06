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

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	novel, err := c.novelService.GetNovelByUpdatesID(idParam)
	if err != nil {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}

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
	statusUpdates := make(chan dtos.ChapterStatus) // New channel for status updates
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

	// Initialize chapter statuses
	chapterStatuses := make(map[any]string, latestChapter+2)
	for i := 1; i <= latestChapter; i++ {
		chapterStatuses[i] = "to download"
	}

	// Send initial status
	fmt.Fprintf(ctx.Writer, "event: status\ndata: %s\n\n", utils.GetStatusJSON(chapterStatuses))
	ctx.Writer.Flush()

	// Worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for chapterNo := range chapterQueue {
				// Notify main loop that we're downloading this chapter
				statusUpdates <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "downloading"}
				log.Printf("Now importing chapter: %d", chapterNo)

				// Check if chapter already exists
				isChapterCreated := c.chapterService.IsChapterCreated(uint(chapterNo), novelID)
				if isChapterCreated {
					skippedChan <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "skipped"}
					continue
				}

				// Import the chapter
				result, err := c.chapterService.ImportChapter(novelUpdatesID, chapterNo)
				if err != nil {
					errorsChan <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "error", Message: err.Error()}
					continue
				}

				results <- result
			}
		}()
	}

	// Populate queue with chapters to process
	go func() {
		defer close(chapterQueue)
		for chapterNo := 1; chapterNo <= latestChapter; chapterNo++ {
			select {
			case <-done:
				return
			case chapterQueue <- chapterNo:
				// Notify main loop that this chapter is queued
				statusUpdates <- dtos.ChapterStatus{ChapterNo: chapterNo, Status: "in queue"}
				log.Printf("Added chapter %d to queue", chapterNo)
			}
		}
	}()

	// Close channels when all workers finish
	go func() {
		wg.Wait()
		close(results)
		close(errorsChan)
		close(skippedChan)
		close(statusUpdates)
	}()

	// Main processing loop (single-threaded map updates)
	for {
		select {
		case update := <-statusUpdates:
			chapterStatuses[update.ChapterNo] = update.Status
		case err := <-errorsChan:
			chapterStatuses[err.ChapterNo] = err.Status
			log.Printf("Error importing chapter %d: %s", err.ChapterNo, err.Status)
		case skipped := <-skippedChan:
			chapterCount++
			chapterStatuses[skipped.ChapterNo] = skipped.Status
			log.Printf("Chapter %d skipped (already exists)", skipped.ChapterNo)
		case result, ok := <-results:
			if !ok {
				// All chapters processed
				log.Printf("All chapters processed (total: %d)", chapterCount-2)
				fmt.Fprintf(ctx.Writer, "event: complete\ndata: All %d chapters processed\n\n", chapterCount)
				ctx.Writer.Flush()
				return nil
			}

			// Save the chapter
			if err := c.chapterService.CreateChapter(novelID, result); err != nil {
				chapterStatuses[int(result.ID)] = "save error: " + err.Error()
				log.Printf("Failed to save chapter %d: %v", result.ID, err)
			} else {
				chapterCount++
				chapterStatuses[int(result.ID)] = "downloaded"
				log.Printf("Successfully processed chapter %d (%d/%d)", result.ID, chapterCount, latestChapter)
			}
		}

		// Send updated status (single-threaded, no race condition)
		fmt.Fprintf(ctx.Writer, "event: status\ndata: %s\n\n", utils.GetStatusJSON(chapterStatuses))
		ctx.Writer.Flush()
	}
}

func (c *ChapterController) GetChapterByNovelUpdatesIDAndChapterNo(ctx *gin.Context) {
	novelTitle := ctx.Param("novel_title")
	chapterNo := ctx.Param("chapter_no")

	chapterNoUint, err := utils.ParseUintID(chapterNo)
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
