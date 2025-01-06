package models

import "gorm.io/gorm"

type Volume struct {
	gorm.Model
	Title        string `gorm:"size:255;not null" json:"title"`
	StartChapter uint   `gorm:"default:0" json:"start_chapter"`
	EndChapter   uint   `gorm:"default:0" json:"final_chapter"`
}