package dtos

type ErrorResponse struct {
	Code  string `json:"code"`  // Error code
	Error string `json:"error"` // Error message
}
