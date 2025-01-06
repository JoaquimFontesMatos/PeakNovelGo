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
	Volumes         []Volume `json:"volumes"`
	Genres          []string `json:"genres"`
}

func ConvertToNovel(imported ImportedNovel) Novel {
	// Convert Tags from []string to []Tag
	tags := make([]Tag, len(imported.Tags))
	for i, tagName := range imported.Tags {
		tags[i] = Tag{Name: tagName}
	}

	// Convert Authors from []string to []Author
	authors := make([]Author, len(imported.Authors))
	for i, authorName := range imported.Authors {
		authors[i] = Author{Name: authorName}
	}

	// Convert Genres from []string to []Genre
	genres := make([]Genre, len(imported.Genres))
	for i, genreName := range imported.Genres {
		genres[i] = Genre{Name: genreName}
	}

	// Return the Novel struct
	return Novel{
		Url:             imported.Url,
		Title:           imported.Title,
		Synopsis:        imported.Synopsis,
		CoverUrl:        imported.CoverUrl,
		Language:        imported.Language,
		Status:          imported.Status,
		NovelUpdatesUrl: imported.NovelUpdatesUrl,
		Tags:            tags,
		Authors:         authors,
		Volumes:         imported.Volumes,
		Genres:          genres,
	}
}

