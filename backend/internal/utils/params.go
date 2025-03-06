package utils

import (
	"backend/internal/types/errors"
	"strconv"
)

const (
	DefaultPage  = 1
	MinPage      = 1
	MaxPage      = 1000
	DefaultLimit = 10
	MinLimit     = 10
	MaxLimit     = 100
)

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
