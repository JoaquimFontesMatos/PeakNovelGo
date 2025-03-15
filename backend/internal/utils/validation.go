package utils

// The function `IsPageOutOfRange` checks if a given page number exceeds the total number of pages
// based on the total items and limit per page.
//
// Parameters:
//   - page (The current page number)
//   - total (The total number of items)
//   - limit (The limit per page)
//
// Returns:
//   - bool (True if the page number exceeds the total number of pages, False otherwise)
func IsPageOutOfRange(page int, total int64, limit int) bool {
	totalPages := CalculateTotalPages(total, limit)
	return int64(page) > totalPages
}
