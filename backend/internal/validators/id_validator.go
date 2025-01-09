package validators

import (
	"errors"
	"strconv"
)

func ValidateID(idParam string) (uint, error) {
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		return 0, errors.New("invalid ID")
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	return uid, nil
}
