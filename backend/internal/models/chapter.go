package models

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	ChapterNo  uint   `gorm:"index,default:0" json:"chapterNo"`
	NovelID    *uint  `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	Title      string `gorm:"size:255;not null" json:"title"`
	ChapterUrl string `gorm:"size:255;not null,uniqueIndex" json:"chapterUrl"`
	Body       string `gorm:"not null" json:"body"`
}

type ImportedChapter struct {
	ID         uint   `json:"id"`
	NovelID    *uint  `json:"novel_id"`
	Title      string `json:"title"`
	ChapterUrl string `json:"url"`
	Body       string `json:"body"`
}

type ImportedChapterMetadata struct {
	Title      string `json:"title" binding:"required"`
	ChapterUrl string `json:"url" binding:"required"`
	Body       string `json:"body" binding:"required"`
	ID         uint   `json:"id"`
}


func (c *ImportedChapter) ToChapter() *Chapter {
	return &Chapter{
		ChapterNo:  c.ID,
		NovelID:    c.NovelID,
		Title:      c.Title,
		ChapterUrl: c.ChapterUrl,
		Body:       c.Body,
	}
}
