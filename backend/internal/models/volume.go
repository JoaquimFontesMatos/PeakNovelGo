package models

import "gorm.io/gorm"

type Volume struct {
	gorm.Model
	NovelID      *uint  `gorm:"constraint:OnDelete:SET NULL;" json:"novel_id"`
	Title        string `gorm:"size:255;not null" json:"title"`
	StartChapter uint   `gorm:"default:0" json:"start_chapter"`
	EndChapter   uint   `gorm:"default:0" json:"end_chapter"`
}