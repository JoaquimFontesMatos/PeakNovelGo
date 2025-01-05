package models

type NovelAuthor struct {
	NovelID  *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	AuthorID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"author_id"`
}