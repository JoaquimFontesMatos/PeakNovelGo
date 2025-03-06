package utils

func IsPageOutOfRange(page int, total int64, limit int) bool {
	totalPages := CalculateTotalPages(total, limit)
	return int64(page) > totalPages
}