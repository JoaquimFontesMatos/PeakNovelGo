package utils

import (
	"github.com/gin-gonic/gin"
)

func BuildPaginatedResponse(c *gin.Context, data interface{}, total int64, page int, limit int) {
	totalPages := CalculateTotalPages(total, limit)

	c.JSON(200, gin.H{
		"data":       data,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

func CalculateTotalPages(total int64, limit int) int64 {
	if total == 0 {
		return 0
	}
	return (total + int64(limit) - 1) / int64(limit)
}