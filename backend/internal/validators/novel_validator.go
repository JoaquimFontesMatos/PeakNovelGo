package validators

import "backend/internal/types"

func ValidateAuthor(author string) error {
	if author == "" {
		return types.WrapError(types.VALIDATION_ERROR, "Author is required", nil)
	}
	if author == " " {
		return types.WrapError(types.VALIDATION_ERROR, "Author is required", nil)
	}
	if len(author) < 1 {
		return types.WrapError(types.VALIDATION_ERROR, "Author must be at least 1 characters long", nil)
	}
	if len(author) > 255 {
		return types.WrapError(types.VALIDATION_ERROR, "Author cannot be longer than 255 characters", nil)
	}
	return nil
}

func ValidateTag(tag string) error {
	if tag == "" {
		return types.WrapError(types.VALIDATION_ERROR, "Tag is required", nil)
	}
	if tag == " " {
		return types.WrapError(types.VALIDATION_ERROR, "Tag is required", nil)
	}
	if len(tag) < 1 {
		return types.WrapError(types.VALIDATION_ERROR, "Tag must be at least 1 characters long", nil)
	}
	if len(tag) > 255 {
		return types.WrapError(types.VALIDATION_ERROR, "Tag cannot be longer than 255 characters", nil)
	}
	return nil
}

func ValidateGenre(genre string) error {
	if genre == "" {
		return types.WrapError(types.VALIDATION_ERROR, "Genre is required", nil)
	}
	if genre == " " {
		return types.WrapError(types.VALIDATION_ERROR, "Genre is required", nil)
	}
	if len(genre) < 1 {
		return types.WrapError(types.VALIDATION_ERROR, "Genre must be at least 1 characters long", nil)
	}
	if len(genre) > 255 {
		return types.WrapError(types.VALIDATION_ERROR, "Genre cannot be longer than 255 characters", nil)
	}
	return nil
}
