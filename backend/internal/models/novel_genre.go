package models

type NovelGenre struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	GenreID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"genreId"`
}