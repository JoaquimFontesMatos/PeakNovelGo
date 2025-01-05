package models

type NovelTag struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novel_id"`
	TagID   *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"tag_id"`
}