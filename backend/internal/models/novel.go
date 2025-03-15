package models

import (
	"gorm.io/gorm"
)

// Novel represents a novel in the database.
//
// Fields:
//   - Title (string): The title of the novel.  Must be unique.  Maximum length 200 characters.
//   - Synopsis (string): A synopsis of the novel.  Required field. Maximum length 5000 characters.
//   - CoverUrl (string): The URL of the novel's cover image. Required field. Maximum length 255 characters.
//   - Language (string): The language of the novel. Required field. Maximum length 255 characters.
//   - Status (string): The status of the novel (e.g., "Ongoing", "Completed"). Required field. Maximum length 500 characters.
//   - NovelUpdatesUrl (string): The URL of the novel on NovelUpdates. Maximum length 255 characters.
//   - NovelUpdatesID (string): The NovelUpdates ID of the novel. Must be unique. Maximum length 255 characters.
//   - Tags ([]Tag): A slice of tags associated with the novel.
//   - Authors ([]Author): A slice of authors of the novel.
//   - Genres ([]Genre): A slice of genres the novel belongs to.
//   - Year (string): The year the novel was published or first released. Required field.
//   - ReleaseFrequency (string): How often new chapters are released (e.g., "Weekly", "Daily"). Required field. Maximum length 255 characters.
//   - LatestChapter (int): The number of the latest chapter. Required field.
type Novel struct {
	gorm.Model
	Title            string   `gorm:"size:200;uniqueIndex" json:"title"`
	Synopsis         string   `gorm:"size:25000;not null" json:"synopsis"`
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
	LatestChapter    int      `gorm:"not null" json:"latestChapter"`
}

// ImportedNovel represents a novel imported from an external source.
//
// Fields:
//   - Title (string): The title of the novel.
//   - Synopsis (string): A brief description of the novel.
//   - CoverUrl (string): URL to the novel's cover image.
//   - Language (ImportedLanguage): The language of the novel.
//   - Status (string): The current status of the novel (e.g., "Ongoing", "Completed").
//   - Tags (Tag): A slice of tags associated with the novel.
//   - Authors (Author): A slice of authors of the novel.
//   - Genres (Genre): A slice of genres the novel belongs to.
//   - Year (string): The year the novel was published or first released.
//   - ReleaseFrequency (string): How often new chapters are released (e.g., "Daily", "Weekly").
//   - LatestChapter (string): The title or identifier of the latest chapter.
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

// ImportedLanguage represents a language imported into a novel.
//
// Parameters:
//   - Name (string): The name of the imported language.
type ImportedLanguage struct {
	Name string `json:"name"`
}
