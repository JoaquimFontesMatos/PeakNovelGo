package models

import (
	"gorm.io/gorm"
)

type Novel struct {
	gorm.Model
	Url             string   `gorm:"size:255" json:"url"`
	Title           string   `gorm:"size:200;uniqueIndex" json:"title"`
	Synopsis        string   `gorm:"size:2000;not null" json:"synopsis"`
	CoverUrl        string   `gorm:"size:255;not null" json:"cover_url"`
	Language        string   `gorm:"size:255;not null" json:"language"`
	Status          string   `gorm:"size:255;not null" json:"status"`
	NovelUpdatesUrl string   `gorm:"size:255" json:"novel_updates_url"`
	Tags            []Tag    `gorm:"many2many:novel_tags;" json:"tags"`
	Authors         []Author `gorm:"many2many:novel_authors;" json:"authors"`
	Volumes         []Volume `gorm:"many2many:novel_volumes;" json:"volumes"`
}

type Tag struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:255;unique;not null" json:"name"`
}

type NovelTag struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	TagID   *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"tag_id"`
}

type Chapter struct {
	gorm.Model
	NovelID    *uint  `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	VolumeID   *uint  `gorm:"constraint:OnDelete:SET NULL;index" json:"volume_id"`
	Title      string `gorm:"size:255;not null" json:"title"`
	ChapterUrl string `gorm:"size:255;not null" json:"chapter_url"`
	Body       string `gorm:"not null" json:"body"`
}

type Volume struct {
	gorm.Model
	NovelID      *uint  `gorm:"constraint:OnDelete:SET NULL;" json:"novel_id"`
	Title        string `gorm:"size:255;not null" json:"title"`
	StartChapter uint   `gorm:"default:0" json:"start_chapter"`
	EndChapter   uint   `gorm:"default:0" json:"end_chapter"`
}

type Author struct {
	gorm.Model
	Name string `gorm:"size:255;unique;not null" json:"name"`
}

type NovelAuthor struct {
	NovelID  *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	AuthorID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"author_id"`
}
