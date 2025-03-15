package models

// NovelGenre represents a many-to-many relationship between novels and genres.
//
// Fields:
//   - NovelID (*uint): The ID of the novel.  Can be NULL.
//   - GenreID (*uint): The ID of the genre. Can be NULL.
type NovelGenre struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	GenreID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"genreId"`
}
