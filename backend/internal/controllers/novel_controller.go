package controllers

import (
	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/types/errors"
	"backend/internal/utils"
	"fmt"
	"sync"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// NovelController struct manages novel-related operations.
//
// Fields:
//   - novelService (interfaces.NovelServiceInterface): An interface that provides access to novel data and operations.
type NovelController struct {
	novelService interfaces.NovelServiceInterface
}

// NewNovelController creates a new NovelController instance.
//
// Parameters:
//   - novelService (interfaces.NovelServiceInterface): The novel service to be used by the controller.
//
// Returns:
//   - *NovelController: A pointer to the newly created NovelController.
func NewNovelController(novelService interfaces.NovelServiceInterface) *NovelController {
	return &NovelController{novelService: novelService}
}

// HandleImportNovelByNovelUpdatesID handles the import of a novel using its NovelUpdates ID.
//
// This function retrieves a novel's data from NovelUpdates using its ID, creates a new novel entry in the database, and
// returns the created novel.
// If an error occurs during the process, it logs the error and returns an appropriate HTTP error response.
//
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - nil:  If the novel was successfully created and the response was sent.
//
// Error types:
//   - error:  Any error that occurs during the creation of the novel.  This will result in an appropriate HTTP error
//     response being sent.

// HandleImportNovelByNovelUpdatesID handles the import of a novel using its NovelUpdates ID.
//
// @Summary Get novel by ID
// @Description This function retrieves a novel's data from NovelUpdates using its ID, creates a new novel entry in the database, and returns the created novel. If an error occurs during the process, it logs the error and returns an appropriate HTTP error response.
// @Tags Novels
// @Accept json
// @Produce json
// @Param novel_updates_id path string true "NovelUpdatesID"
// @Success 200 {object} models.Novel
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Security BearerAuth
// @Router /novels/{novel_updates_id} [post]
func (n *NovelController) HandleImportNovelByNovelUpdatesID(ctx *gin.Context) {
	novelUpdatesID := ctx.Param("novel_updates_id")

	createdNovel, err := n.novelService.CreateNovel(novelUpdatesID)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(http.StatusCreated, createdNovel)
}

// GetNovelsByAuthorName retrieves novels by author name, handling pagination and validation.
//
// @Summary Get novels by author name
// @Description Retrieves a paginated list of novels written by a specific author.
// @Tags Novels
// @Accept json
// @Produce json
// @Param author_name path string true "Author name"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} dtos.PaginatedResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/authors/{author_name} [get]
func (n *NovelController) GetNovelsByAuthorName(ctx *gin.Context) {
	authorName := ctx.Param("author_name")

	// Parse parameters
	page, err := utils.ParsePage(ctx.Query("page"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	limit, err := utils.ParseLimit(ctx.Query("limit"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	novels, total, err := n.novelService.GetNovelsByAuthorName(authorName, page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Validate results
	if total == 0 {
		utils.HandleError(ctx, errors.ErrNoResults)
		return
	}

	if utils.IsPageOutOfRange(page, total, limit) {
		utils.HandleError(ctx, errors.ErrPageOutOfRange)
		return
	}

	// Build response
	utils.BuildPaginatedResponse(ctx, novels, total, page, limit)
}

// GetNovelsByGenreName retrieves novels by genre name, handling pagination and validation.
//
// @Summary Get novels by genre name
// @Description Retrieves a paginated list of novels written by a specific genre.
// @Tags Novels
// @Accept json
// @Produce json
// @Param genre_name path string true "Genre name"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} dtos.PaginatedResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/genres/{genre_name} [get]
func (n *NovelController) GetNovelsByGenreName(ctx *gin.Context) {
	genreName := ctx.Param("genre_name")

	// Parse parameters
	page, err := utils.ParsePage(ctx.Query("page"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	limit, err := utils.ParseLimit(ctx.Query("limit"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	novels, total, err := n.novelService.GetNovelsByGenreName(genreName, page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Validate results
	if total == 0 {
		utils.HandleError(ctx, errors.ErrNoResults)
		return
	}

	if utils.IsPageOutOfRange(page, total, limit) {
		utils.HandleError(ctx, errors.ErrPageOutOfRange)
		return
	}

	// Build response
	utils.BuildPaginatedResponse(ctx, novels, total, page, limit)
}

// GetNovelsByTagName retrieves novels by tag name, handling pagination and validation.
//
// @Summary Get novels by tag name
// @Description Retrieves a paginated list of novels written by a specific tag.
// @Tags Novels
// @Accept json
// @Produce json
// @Param tag_name path string true "Tag name"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} dtos.PaginatedResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/tags/{tag_name} [get]
func (n *NovelController) GetNovelsByTagName(ctx *gin.Context) {
	tagName := ctx.Param("tag_name")

	// Parse parameters
	page, err := utils.ParsePage(ctx.Query("page"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	limit, err := utils.ParseLimit(ctx.Query("limit"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	novels, total, err := n.novelService.GetNovelsByTagName(tagName, page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Validate results
	if total == 0 {
		utils.HandleError(ctx, errors.ErrNoResults)
		return
	}

	if utils.IsPageOutOfRange(page, total, limit) {
		utils.HandleError(ctx, errors.ErrPageOutOfRange)
		return
	}

	// Build response
	utils.BuildPaginatedResponse(ctx, novels, total, page, limit)
}

// GetNovels retrieves a paginated list of novels.
//
// @Summary Get novels
// @Description Retrieves a paginated list of novels.
// @Tags Novels
// @Accept json
// @Produce json
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of items per page (default: 10)"
// @Success 200 {object} dtos.PaginatedResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/ [get]
func (n *NovelController) GetNovels(ctx *gin.Context) {
	// Parse parameters
	page, err := utils.ParsePage(ctx.Query("page"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	limit, err := utils.ParseLimit(ctx.Query("limit"))
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Get data
	novels, total, err := n.novelService.GetNovels(page, limit)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	// Validate results
	if total == 0 {
		utils.HandleError(ctx, errors.ErrNoResults)
		return
	}

	if utils.IsPageOutOfRange(page, total, limit) {
		utils.HandleError(ctx, errors.ErrPageOutOfRange)
		return
	}

	// Build response
	utils.BuildPaginatedResponse(ctx, novels, total, page, limit)
}

// GetNovelByID retrieves a novel by its ID.
//
// @Summary Get novel by ID
// @Description Retrieves a novel by its ID.
// @Tags Novels
// @Accept json
// @Produce json
// @Param novel_id path string true "Novel ID"
// @Success 200 {object} models.Novel
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/{novel_id} [get]
func (n *NovelController) GetNovelByID(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := utils.ParseUintID(idParam)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	novel, err := n.novelService.GetNovelByID(id)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

// GetNovelByUpdatesID retrieves a novel based on its title (which acts as the updates ID).
//
// @Summary Get novel by NovelByUpdatesID
// @Description Retrieves a novel based on its title (which acts as the updates ID).
// @Tags Novels
// @Accept json
// @Produce json
// @Param title path string true "NovelUpdatesID"
// @Success 200 {object} models.Novel
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Failure 503 {object} dtos.ErrorResponse
// @Router /novels/title/{title} [get]
func (n *NovelController) GetNovelByUpdatesID(ctx *gin.Context) {
	title := ctx.Param("title")

	novel, err := n.novelService.GetNovelByUpdatesID(title)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

// HandleBatchUpdateNovels handles streaming response for updating novels
func (n *NovelController) HandleBatchUpdateNovels(ctx *gin.Context) {
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Set("Connection", "keep-alive")
	ctx.Writer.Flush()

	err := n.processNovelsWithStreaming(ctx)
	if err != nil {
		fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", err.Error())
		ctx.Writer.Flush()
		return
	}
}

func (n *NovelController) processNovelsWithStreaming(ctx *gin.Context) error {
	const workerCount = 2
	novelCount := 0
	novelQueue := make(chan int, workerCount)
	results := make(chan string)
	errorsChan := make(chan dtos.NovelStatus)
	statusUpdates := make(chan dtos.NovelStatus) // New channel for status updates
	done := make(chan struct{})

	var wg sync.WaitGroup

	// Fetch all novels
	novels, total, err := n.novelService.GetNovels(1, 999999999)
	if err != nil {
		sendSSEError(ctx, "Failed to fetch novels: "+err.Error())
		return err
	}

	totalInt := int(total)
	if totalInt == 0 {
		sendSSEError(ctx, "No novels found")
		return fmt.Errorf("no novels found")
	}

	// Initialize novel statuses (thread-safe updates via main loop)
	novelStatuses := make(map[any]string, totalInt)
	for i := 0; i < totalInt; i++ {
		id := novels[i].NovelUpdatesID
		novelStatuses[id] = "to update"
	}
	sendSSEStatus(ctx, novelStatuses)

	// Worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for novelNo := range novelQueue {
				id := novels[novelNo].NovelUpdatesID

				// Notify main loop of status change
				statusUpdates <- dtos.NovelStatus{NovelUpdatesId: id, Status: "updating"}

				// Process the novel
				_, err := n.novelService.CreateNovel(id)
				if err != nil {
					errorsChan <- dtos.NovelStatus{NovelUpdatesId: id, Status: "error", Message: err.Error()}
					continue
				}

				results <- id
			}
		}()
	}

	// Populate queue with novels to process
	go func() {
		defer close(novelQueue)
		for i := 0; i < totalInt; i++ {
			select {
			case <-done:
				return
			case novelQueue <- i:
				id := novels[i].NovelUpdatesID
				statusUpdates <- dtos.NovelStatus{NovelUpdatesId: id, Status: "in queue"}
			}
		}
	}()

	// Close channels when workers finish
	go func() {
		wg.Wait()
		close(results)
		close(errorsChan)
		close(statusUpdates)
	}()

	// Main loop (single-threaded map updates)
	for {
		select {
		case update := <-statusUpdates:
			novelStatuses[update.NovelUpdatesId] = update.Status
			sendSSEStatus(ctx, novelStatuses)

		case err := <-errorsChan:
			log.Printf("error: %s", err.Message)
			novelStatuses[err.NovelUpdatesId] = err.Status
			sendSSEStatus(ctx, novelStatuses)

		case result, ok := <-results:
			if !ok {
				sendSSEComplete(ctx, fmt.Sprintf("All %d novels processed", novelCount))
				return nil
			}
			novelCount++
			novelStatuses[result] = "updated"
			sendSSEStatus(ctx, novelStatuses)
		}
	}
}

// Helper functions for SSE messages
func sendSSEStatus(ctx *gin.Context, statuses map[any]string) {
	fmt.Fprintf(ctx.Writer, "event: status\ndata: %s\n\n", utils.GetStatusJSON(statuses))
	ctx.Writer.Flush()
}

func sendSSEError(ctx *gin.Context, message string) {
	fmt.Fprintf(ctx.Writer, "event: error\ndata: %s\n\n", message)
	ctx.Writer.Flush()
}

func sendSSEComplete(ctx *gin.Context, message string) {
	fmt.Fprintf(ctx.Writer, "event: complete\ndata: %s\n\n", message)
	ctx.Writer.Flush()
}
