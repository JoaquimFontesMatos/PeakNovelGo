package utils

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"backend/internal/types"
	"backend/internal/types/errors"
)

// The ParseUintID function takes a string parameter representing an ID, converts it to a uint, and returns
// it along with any errors encountered.
//
// Parameters:
//   - idParam string (string representing the ID)
//
// Returns:
//   - uint (parsed ID)
//   - error (errors.ErrParseUint if the string cannot be parsed to a uint)
func ParseUintID(idParam string) (uint, error) {
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		return 0, types.WrapError(errors.PARSE_UINT_ERROR, "Invalid ID", http.StatusBadRequest, err)
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	return uid, nil
}

// The ParseInt function in Go parses a string to an integer, handling whitespace and newline
// characters, and returns the integer value or an error.
//
// Parameters:
//   - int string (string to parse)
//
// Returns:
//   - int (parsed integer)
//   - error (errors.ErrParseInt if the string cannot be parsed to an integer)
func ParseInt(int string) (int, error) {
	int = strings.TrimSpace(int)
	int = strings.ReplaceAll(int, "\n", "")

	num, err := strconv.Atoi(int)

	if err != nil || num <= 0 {
		return 0, types.WrapError(errors.PARSE_INT_ERROR, "Invalid Conversion", http.StatusBadRequest, err)
	}

	return num, nil
}

// NovelUpdatesIDParser extracts the NovelUpdates novel ID from a given string.
type NovelUpdatesIDParser struct {
	// regex is the compiled regular expression used for ID extraction.
	regex *regexp.Regexp
}

// NewNovelUpdatesIDParser creates a new NovelUpdatesIDParser instance.
//
// Returns:
//   - *NovelUpdatesIDParser (pointer to NovelUpdatesIDParser struct)
func NewNovelUpdatesIDParser() *NovelUpdatesIDParser {
	return &NovelUpdatesIDParser{
		regex: regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`), // Updated regex
	}
}

// Parse parses a NovelUpdates ID, validating its length and format.
//
// Parameters:
//   - novelUpdatesID (string): The NovelUpdates ID to parse.
//
// Returns:
//   - string: The parsed and validated NovelUpdates ID.  Returns an empty string if the ID is invalid.
//   - error: An error if the NovelUpdates ID is invalid.  Possible errors include `errors.ErrInvalidNovelUpdatesID`.
//
// Error types:
//   - errors.ErrInvalidNovelUpdatesID: Returned if the provided ID is invalid (e.g., incorrect length or format).
func (n *NovelUpdatesIDParser) Parse(novelUpdatesID string) (string, error) {
	if len(novelUpdatesID) < 1 || len(novelUpdatesID) > 255 {
		return "", errors.ErrInvalidNovelUpdatesID
	}

	// Convert to lowercase and replace spaces with dashes
	lowercaseNovelUpdatesID := strings.ToLower(novelUpdatesID)
	noSpacesNovelUpdatesID := strings.ReplaceAll(lowercaseNovelUpdatesID, " ", "-")

	// Validate the ID using the regex
	if !n.regex.MatchString(noSpacesNovelUpdatesID) {
		return "", errors.ErrInvalidNovelUpdatesID
	}

	return noSpacesNovelUpdatesID, nil
}
