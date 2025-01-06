package utils

import (
	"regexp"
	"html"
)

// StripHTML removes HTML tags while preserving paragraph breaks and decodes HTML entities.
func StripHTML(input string) string {
	// Replace <p> and </p> tags with newlines
	reParagraph := regexp.MustCompile(`(?i)<\/?p\s*>`)
	input = reParagraph.ReplaceAllString(input, "\n")

	// Remove all other HTML tags
	reTags := regexp.MustCompile(`<.*?>`)
	cleaned := reTags.ReplaceAllString(input, "")

	// Decode HTML entities (e.g., &nbsp; -> space)
	return html.UnescapeString(cleaned)
}
