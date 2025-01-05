package models

type NovelGenre struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	GenreID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"genre_id"`
}