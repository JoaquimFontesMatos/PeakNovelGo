package models

type NovelAuthor struct {
	NovelID  *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	AuthorID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"authorId"`
}