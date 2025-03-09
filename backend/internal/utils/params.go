package utils

import (
	"backend/internal/types/errors"
	"strconv"
)

const (
	// DefaultPage defines the default page value for operations when no specific page is provided.
	DefaultPage  = 1
	// MinPage defines the minimum allowable page value for a given operation, ensuring it is not less than the defined value.
	MinPage      = 1
	// MaxPage defines the maximum allowable page value for a given operation, ensuring it does not exceed the defined value.
	MaxPage      = 1000

	// DefaultLimit defines the default limit value for operations when no specific limit is provided.
	DefaultLimit = 10
	// MinLimit defines the minimum allowable limit value for a given operation, ensuring it is not less than the defined value.
	MinLimit = 10
	// MaxLimit defines the maximum allowable limit value for a given operation, ensuring it does not exceed the defined value.
	MaxLimit = 100
)

// ParsePage parses the given page string into an integer, validating it against minimum and maximum allowable values.
//
// Parameters:
//	- pageStr string (page string to be parsed)
//
// Returns:
//	- int (the default page if the input is empty)
//	- error (errors.ErrInvalidPage if the input is invalid or out of range)
func ParsePage(pageStr string) (int, error) {
	if pageStr == "" {
		return DefaultPage, nil
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < MinPage || page > MaxPage {
		return 0, errors.ErrInvalidPage
	}
	return page, nil
}

// ParseLimit parses the given limit string into an integer, validating it against minimum and maximum allowable values.
//
// Parameters:
//	- limitStr string (limit string to be parsed)
//
// Returns:
//	- int (the default limit if the input is empty)
//	- error (errors.ErrInvalidLimit if the input is invalid or out of range)
func ParseLimit(limitStr string) (int, error) {
	if limitStr == "" {
		return DefaultLimit, nil
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < MinLimit || limit > MaxLimit {
		return 0, errors.ErrInvalidLimit
	}
	return limit, nil
}
