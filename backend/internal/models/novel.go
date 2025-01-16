package models

import (
	"gorm.io/gorm"
)

type Novel struct {
	gorm.Model
	Url             string   `gorm:"size:255" json:"url"`
	Title           string   `gorm:"size:200;uniqueIndex" json:"title"`
	Synopsis        string   `gorm:"size:2000;not null" json:"synopsis"`
	CoverUrl        string   `gorm:"size:255;not null" json:"coverUrl"`
	Language        string   `gorm:"size:255;not null" json:"language"`
	Status          string   `gorm:"size:255;not null" json:"status"`
	NovelUpdatesUrl string   `gorm:"size:255" json:"novel_updatesUrl"`
	Tags            []Tag    `gorm:"many2many:novel_tags;" json:"tags"`
	Authors         []Author `gorm:"many2many:novel_authors;" json:"authors"`
	Genres          []Genre  `gorm:"many2many:novel_genres;" json:"genres"`
}

type ImportedNovel struct {
	Url             string   `json:"url"`
	Title           string   `json:"title"`
	Synopsis        string   `json:"synopsis"`
	CoverUrl        string   `json:"cover_url"`
	Language        string   `json:"language"`
	Status          string   `json:"status"`
	NovelUpdatesUrl string   `json:"novel_updates_url"`
	Tags            []string `json:"novel_tags"`
	Authors         []string `json:"authors"`
	Genres          []string `json:"genres"`
}