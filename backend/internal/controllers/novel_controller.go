package controllers

import (
	"backend/internal/services/interfaces"
	"backend/internal/types/errors"
	"backend/internal/utils"

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
// Parameters:
//   - ctx (*gin.Context): Gin context.
//
// Returns:
//   - :  Paginated list of novels is sent as a JSON response. An error response is sent if any error occurs.
//
// Error types:
//   - errors.ErrNoResults: Returned when no novels are found for the given author.
//   - errors.ErrPageOutOfRange: Returned when the requested page is out of range.
//   - error:  Various other errors during database interaction or parameter parsing.
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

// GetNovelsByGenreName retrieves novels based on genre name, pagination, and limit.
//
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - :  Paginated list of novels is sent as a JSON response. An error response is sent if any error occurs.
//
// Error types:
//   - errors.ErrNoResults: Returned when no novels are found for the given genre.
//   - errors.ErrPageOutOfRange: Returned when the requested page is out of range.
//   - error:  Various other errors during database interaction or parameter parsing.
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

// GetNovelsByTagName retrieves novels based on tag name, pagination, and limit.
//
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - :  Paginated list of novels is sent as a JSON response. An error response is sent if any error occurs.
//
// Error types:
//   - errors.ErrNoResults: Returned when no novels are found for the given tag.
//   - errors.ErrPageOutOfRange: Returned when the requested page is out of range.
//   - error:  Various other errors during database interaction or parameter parsing.
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
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - : a paginated list of novels. Returns a JSON response containing the novels, total count, page number, and limit.
//     Error responses are handled within the function.
//
// Error types:
//   - *errors.Error: Various error types are handled, including database errors, invalid parameters, and out-of-range
//     page requests. Specific error types are defined in the `errors` package.
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
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - : a novel. Returns a JSON response containing the fetched novel.
//     Error responses are handled within the function.
//
// Error types:
//   - *utils.Error: if the novel is not found or another error occurs during retrieval.  The specific error will be
//     handled and responded to via `utils.HandleError`.
//   - error: if the provided novel ID is invalid.  A 400 Bad Request response will be returned with the error message.
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
// Parameters:
//   - ctx (*gin.Context): the Gin context.
//
// Returns:
//   - : a novel. Returns a JSON response containing the fetched novel.
//     Error responses are handled within the function.
//
// Error types:
//   - error: various errors may occur during database retrieval or other operations.  The specific error will be handled
//     and returned via the Gin context.
func (n *NovelController) GetNovelByUpdatesID(ctx *gin.Context) {
	title := ctx.Param("title")

	novel, err := n.novelService.GetNovelByUpdatesID(title)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, novel)
}
