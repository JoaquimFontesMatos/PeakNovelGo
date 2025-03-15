package dtos

type PaginatedResponse struct {
	Data       interface{} `json:"data"`  // The paginated data (e.g., list of novels)
	Total      int64       `json:"total"` // Total number of items
	Page       int         `json:"page"`  // Current page number
	Limit      int         `json:"limit"` // Number of items per page
	TotalPages int64       `json:"totalPages"`
}
