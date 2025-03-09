package models

import (
	"gorm.io/gorm"
)

// BookmarkedNovel represents a user's bookmarked novel with its associated metadata.
//
// Fields:
//   - NovelID (int): The ID of the novel bookmarked. The constraint `OnUpdate:CASCADE,OnDelete:SET NULL;` ensures that
//
// if the novel is updated or deleted, the corresponding entry in this table is updated or set to NULL respectively.
//   - UserID (int): The ID of the user who bookmarked the novel. The constraint `OnUpdate:CASCADE,OnDelete:SET NULL;`
//
// ensures that if the user is updated or deleted, the corresponding entry in this table is updated or set to NULL respectively.
//   - Status (string): The status of the bookmarked novel (e.g., "reading", "completed", "dropped").  This field is
//
// required (`not null`).
//   - Score (int): The user's rating of the novel (0-5). Defaults to 0.
//   - CurrentChapter (int): The chapter the user has currently reached. Defaults to 0.
type BookmarkedNovel struct {
	gorm.Model
	NovelID int `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"novelId"`
	UserID int `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"userId"`
	Status string `gorm:"not null" json:"status"`
	Score int `gorm:"default:0;max:5" json:"score"`
	CurrentChapter int `gorm:"default:0" json:"currentChapter"`
}
