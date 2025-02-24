package models

import (
	"gorm.io/gorm"
)

type Novel struct {
	gorm.Model
	Title            string   `gorm:"size:200;uniqueIndex" json:"title"`
	Synopsis         string   `gorm:"size:5000;not null" json:"synopsis"`
	CoverUrl         string   `gorm:"size:255;not null" json:"coverUrl"`
	Language         string   `gorm:"size:255;not null" json:"language"`
	Status           string   `gorm:"size:500;not null" json:"status"`
	NovelUpdatesUrl  string   `gorm:"size:255" json:"novelUpdatesUrl"`
	NovelUpdatesID   string   `gorm:"size:255;uniqueIndex" json:"novelUpdatesId"`
	Tags             []Tag    `gorm:"many2many:novel_tags;" json:"tags"`
	Authors          []Author `gorm:"many2many:novel_authors;" json:"authors"`
	Genres           []Genre  `gorm:"many2many:novel_genres;" json:"genres"`
	Year             string   `gorm:"not null" json:"year"`
	ReleaseFrequency string   `gorm:"size:255;not null" json:"releaseFrequency"`
	LatestChapter    int     `gorm:"size:255;not null" json:"latestChapter"`
}

type ImportedNovel struct {
	Title            string           `json:"title"`
	Synopsis         string           `json:"description"`
	CoverUrl         string           `json:"image"`
	Language         ImportedLanguage `json:"language"`
	Status           string           `json:"status"`
	Tags             []Tag            `json:"tags"`
	Authors          []Author         `json:"authors"`
	Genres           []Genre          `json:"genre"`
	Year             string           `json:"year"`
	ReleaseFrequency string           `json:"release_freq"`
	LatestChapter    string           `json:"latest_chapter"`
}

type ImportedLanguage struct {
	Name string `json:"name"`
}
