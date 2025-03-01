package controllers

import (
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"backend/internal/types"
	"backend/internal/utils"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookmarkController struct {
	bookmarkService interfaces.BookmarkServiceInterface
}

func NewBookmarkController(bookmarkService interfaces.BookmarkServiceInterface) *BookmarkController {
	return &BookmarkController{bookmarkService: bookmarkService}
}

func (b *BookmarkController) GetBookmarkedNovelsByUserID(ctx *gin.Context) {
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

	novels, total, err := b.bookmarkService.GetBookmarkedNovelsByUserID(id, page, limit)
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

func (b *BookmarkController) CreateBookmark(ctx *gin.Context) {
	var bookmarkedNovel models.BookmarkedNovel
	if err := ctx.ShouldBindJSON(&bookmarkedNovel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdBookmarkedNovel, err := b.bookmarkService.CreateBookmark(bookmarkedNovel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save bookmarked novel to database"})
		return
	}

	log.Println("Bookmarked novel saved successfully")

	ctx.JSON(http.StatusOK, createdBookmarkedNovel)
}

func (b *BookmarkController) UpdateBookmark(ctx *gin.Context) {
	var novel models.BookmarkedNovel
	if err := ctx.ShouldBindJSON(&novel); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedNovel, err := b.bookmarkService.UpdateBookmark(novel)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedNovel)
}

func (b *BookmarkController) GetBookmarkByUserIDAndNovelID(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")
	novelIDParam := ctx.Param("novel_id")

	userID, err := utils.ParseID(userIDParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	novel, err := b.bookmarkService.GetBookmarkByUserIDAndNovelID(userID, novelIDParam)
	if err != nil {
		var myError *types.MyError
		if errors.As(err, &myError) {
			switch myError.Code {
			case types.NOVEL_NOT_FOUND_ERROR:
				ctx.JSON(http.StatusNotFound, gin.H{"error": myError.Message})
			}
			return
		}

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, novel)
}

func (b *BookmarkController) UnbookmarkNovel(ctx *gin.Context) {

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

	err = b.bookmarkService.UnbookmarkNovel(uint(userID), uint(novelID))
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