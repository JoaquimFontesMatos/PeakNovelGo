package controllers

import (
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"errors"

	"log"
	"net/http"
	"strconv"

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

	if novelUpdatesID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No novel updates ID provided"})
		return
	}

	createdNovel, err := n.novelService.CreateNovel(novelUpdatesID)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.CONFLICT_ERROR:
				ctx.JSON(http.StatusConflict, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NOVEL_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			case types.SCRIPT_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			}
			log.Println(err)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(http.StatusCreated, createdNovel)
}

func (n *NovelController) GetNovelsByAuthorName(ctx *gin.Context) {
	authorName := ctx.Param("author_name")

	// Parse and validate page parameter
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	// Parse and validate limit parameter
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 10 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	novels, total, err := n.novelService.GetNovelsByAuthorName(authorName, page, limit)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NO_NOVELS_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if total == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No novels found for the author"})
		return
	}

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

	// Parse and validate page parameter
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	// Parse and validate limit parameter
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 10 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	novels, total, err := n.novelService.GetNovelsByGenreName(genreName, page, limit)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NO_NOVELS_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if total == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No novels found for the genre"})
		return
	}

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

	// Parse and validate page parameter
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	// Parse and validate limit parameter
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 10 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	novels, total, err := n.novelService.GetNovelsByTagName(tagName, page, limit)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NO_NOVELS_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if total == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No novels found for the tag"})
		return
	}

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
	// Parse and validate page parameter
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 || page > 1000 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	// Parse and validate limit parameter
	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 10 || limit > 100 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	novels, total, err := n.novelService.GetNovels(page, limit)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NO_NOVELS_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if total == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No novels found"})
		return
	}

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
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.VALIDATION_ERROR:
				ctx.JSON(http.StatusBadRequest, gin.H{"error": myError.Message})
			case types.DATABASE_ERROR:
				ctx.JSON(http.StatusServiceUnavailable, gin.H{"error": myError.Message})
			case types.NOVEL_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			log.Println(err)
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}
