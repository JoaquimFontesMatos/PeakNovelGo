package controllers

import (
	"backend/internal/services/interfaces"
	"backend/internal/types/errors"
	"backend/internal/utils"

	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NovelController struct {
	novelService interfaces.NovelServiceInterface
}

func NewNovelController(novelService interfaces.NovelServiceInterface) *NovelController {
	return &NovelController{novelService: novelService}
}

// HandleImportNovelByNovelUpdatesID handles POST /novel
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

func (n *NovelController) GetNovelByID(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := utils.ParseID(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novel, err := n.novelService.GetNovelByID(id)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

func (n *NovelController) GetNovelByUpdatesID(ctx *gin.Context) {
	title := ctx.Param("title")

	novel, err := n.novelService.GetNovelByUpdatesID(title)
	if err != nil {
		utils.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, novel)
}
