package models

import "gorm.io/gorm"

type Chapter struct {
	gorm.Model
	NovelID    *uint  `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	Title      string `gorm:"size:255;not null" json:"title"`
	ChapterUrl string `gorm:"size:255;not null,uniqueIndex" json:"chapter_url"`
	Body       string `gorm:"not null" json:"body"`
}
