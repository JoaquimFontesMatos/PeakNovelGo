package utils

import (
	"backend/internal/types"
	"strconv"
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
		return 0, types.WrapError(types.INVALID_ID_ERROR, "Invalid ID", err)
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	return uid, nil
}
