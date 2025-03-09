package utils

import (
	"github.com/gin-gonic/gin"
)

// BuildPaginatedResponse constructs a JSON response with paginated data and sends it to the client using the given context.
// It includes data, total items, current page, items per page, and total pages in the response.
//
// Parameters:
//	- c *gin.Context (context of the response)
//	- data interface{} (paginated data to be sent)
//	- total int64 (total number of entries)
//	- page int (page number)
//	- limit int (limit of entries)
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

// CalculateTotalPages computes the total number of pages based on the total items and items per page (limit).
//
// Parameters:
//	- total int64 (total number of entries)
//	- limit int (limit of entries)
//
// Returns:
//	- int64 (represents the total pages, ensuring proper rounding for any remaining items. If the total is 0, it returns 0)
func CalculateTotalPages(total int64, limit int) int64 {
	if total == 0 {
		return 0
	}
	return (total + int64(limit) - 1) / int64(limit)
}