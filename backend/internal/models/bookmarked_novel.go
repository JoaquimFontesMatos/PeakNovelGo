package models

import (
	"gorm.io/gorm"
)

type BookmarkedNovel struct {
	gorm.Model
	NovelID        int    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"novelId"`
	UserID         int    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"userId"`
	Status         string `gorm:"not null" json:"status"`
	Score          int    `gorm:"default:0;max:5" json:"score"`
	CurrentChapter int    `gorm:"default:0" json:"currentChapter"`
}
