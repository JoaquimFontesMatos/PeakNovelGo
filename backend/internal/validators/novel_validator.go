package validators

import (
	"backend/internal/types/errors"
)

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
