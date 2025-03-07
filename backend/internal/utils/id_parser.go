package utils

import (
	"backend/internal/types"
	"backend/internal/types/errors"
	"net/http"
	"strconv"
	"strings"
)

// ParseID parses the ID from the URL parameter and returns it as a uint.
//
// Parameters:
//   - idParam string (ID parameter from the URL)
//
// Returns:
//   - uint (parsed ID)
//   - INVALID_ID_ERROR if the ID is invalid
func ParseID(idParam string) (uint, error) {
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		return 0, types.WrapError(errors.PARSE_UINT_ERROR, "Invalid ID", http.StatusBadRequest, err)
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	return uid, nil
}

func ParseInt(int string) (int, error) {
	int = strings.TrimSpace(int)
	int = strings.ReplaceAll(int, "\n", "")

	num, err := strconv.Atoi(int)

	if err != nil || num <= 0 {
		return 0, types.WrapError(errors.PARSE_INT_ERROR, "Invalid Conversion", http.StatusBadRequest, err)
	}

	return num, nil
}
