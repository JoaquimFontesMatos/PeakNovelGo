package utils

import (
	"backend/internal/types"
	"regexp"
	"strings"
)

// NovelUpdatesIDParser struct represents a parser for novel updates IDs.
type NovelUpdatesIDParser struct {
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

// Parse parses a novel updates ID.
//
// Parameters:
//   - novelUpdatesID string (novel updates ID)
//
// Returns:
//   - string (novel updates ID)
//   - error (error if the novel updates ID is invalid)
func (n *NovelUpdatesIDParser) Parse(novelUpdatesID string) (string, error) {
	if len(novelUpdatesID) < 1 || len(novelUpdatesID) > 255 {
		return "", types.WrapError(types.VALIDATION_ERROR, "novel updates ID must be 1 to 255 characters long", nil)
	}

	// Convert to lowercase and replace spaces with dashes
	lowercaseNovelUpdatesID := strings.ToLower(novelUpdatesID)
	noSpacesNovelUpdatesID := strings.ReplaceAll(lowercaseNovelUpdatesID, " ", "-")

	// Validate the ID using the regex
	if !n.regex.MatchString(noSpacesNovelUpdatesID) {
		return "", types.WrapError(types.VALIDATION_ERROR, "novel updates ID can only contain lowercase letters, numbers, or single dashes", nil)
	}

	return noSpacesNovelUpdatesID, nil
}