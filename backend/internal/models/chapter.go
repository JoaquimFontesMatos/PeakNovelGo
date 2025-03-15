package models

import "gorm.io/gorm"

// Chapter represents a chapter in a novel.
//
// Fields:
//  - ChapterNo (uint): The chapter number. Defaults to 0.
//  - NovelID (*uint): The ID of the novel this chapter belongs to.  Can be NULL.
//  - Title (string): The title of the chapter (max 255 characters).  Cannot be empty.
//  - ChapterUrl (string): The URL of the chapter (max 255 characters). Must be unique. Cannot be empty.
//  - Body (string): The content of the chapter. Cannot be empty.
type Chapter struct {
	gorm.Model
	ChapterNo  uint   `gorm:"index,default:0" json:"chapterNo"`
	NovelID    *uint  `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	Title      string `gorm:"size:255;not null" json:"title"`
	ChapterUrl string `gorm:"size:255;not null,uniqueIndex" json:"chapterUrl"`
	Body       string `gorm:"not null" json:"body"`
}

// ImportedChapter represents a chapter imported from an external source.
//
// Fields:
//   - ID (uint): The unique identifier of the chapter.
//   - NovelID (*uint): A pointer to the unique identifier of the novel this chapter belongs to.  Can be nil if the novel ID is unknown.
//   - Title (string): The title of the chapter.
//   - ChapterUrl (string): The URL where the chapter was originally sourced from.
//   - Body (string): The content of the chapter.
type ImportedChapter struct {
	ID         uint   `json:"id"`
	NovelID    *uint  `json:"novel_id"`
	Title      string `json:"title"`
	ChapterUrl string `json:"url"`
	Body       string `json:"body"`
}

// ImportedChapterMetadata represents metadata for an imported chapter.
//
// Fields:
//  - Title (string): The title of the chapter.  Required.
//  - ChapterUrl (string): The URL of the chapter. Required.
//  - Body (string): The body content of the chapter. Required.
//  - ID (uint): The unique identifier of the chapter.
type ImportedChapterMetadata struct {
	Title      string `json:"title" binding:"required"`
	ChapterUrl string `json:"url" binding:"required"`
	Body       string `json:"body" binding:"required"`
	ID         uint   `json:"id"`
}

// ToChapter converts an ImportedChapter to a Chapter.
//
// Parameters:
//  - c (*ImportedChapter): The ImportedChapter to convert.
//
// Returns:
//   - *Chapter: A pointer to the converted Chapter.
func (c *ImportedChapter) ToChapter() *Chapter {
	return &Chapter{
		ChapterNo:  c.ID,
		NovelID:    c.NovelID,
		Title:      c.Title,
		ChapterUrl: c.ChapterUrl,
		Body:       c.Body,
	}
}
