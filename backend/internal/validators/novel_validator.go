package validators

import (
	"backend/internal/types/errors"
)

// ValidateAuthor validates the author string.
//
// Parameters:
//   - author (string): The author string to validate.
//
// Returns:
//   - error: nil if the author is valid, otherwise an error indicating the issue.
//
// Error types:
//   - errors.ErrAuthorRequired: if the author string is empty or contains only spaces.
//   - errors.ErrAuthorTooShort: if the author string is less than 1 character long.
//   - errors.ErrAuthorTooLong: if the author string is more than 255 characters long.
func ValidateAuthor(author string) error {
	if author == "" {
		return errors.ErrAuthorRequired
	}
	if author == " " {
		return errors.ErrAuthorRequired
	}
	if len(author) < 1 {
		return errors.ErrAuthorTooShort
	}
	if len(author) > 255 {
		return errors.ErrAuthorTooLong
	}

	return nil
}

// ValidateTag validates a tag string.
//
// Parameters:
//   - tag (string): The tag string to validate.
//
// Returns:
//   - error: nil if the tag is valid, otherwise an error indicating the issue.
//
// Error types:
//   - errors.ErrTagRequired: if the tag string is empty or contains only spaces.
//   - errors.ErrTagTooShort: if the tag string is less than 1 character long.
//   - errors.ErrTagTooLong: if the tag string is more than 255 characters long.
func ValidateTag(tag string) error {
	if tag == "" {
		return errors.ErrTagRequired
	}
	if tag == " " {
		return errors.ErrTagRequired
	}
	if len(tag) < 1 {
		return errors.ErrTagTooShort
	}
	if len(tag) > 255 {
		return errors.ErrTagTooLong
	}

	return nil
}

// ValidateGenre validates the genre string.
//
// Parameters:
//   - genre (string): The genre string to validate.
//
// Returns:
//   - error: nil if the genre is valid, otherwise an error indicating the issue.
//
// Error types:
//   - errors.ErrGenreRequired: if the genre string is empty or contains only spaces.
//   - errors.ErrGenreTooShort: if the genre string is less than 1 character long.
//   - errors.ErrGenreTooLong: if the genre string is more than 255 characters long.
func ValidateGenre(genre string) error {
	if genre == "" {
		return errors.ErrGenreRequired
	}
	if genre == " " {
		return errors.ErrGenreRequired
	}
	if len(genre) < 1 {
		return errors.ErrGenreTooShort
	}
	if len(genre) > 255 {
		return errors.ErrGenreTooLong
	}

	return nil
}
