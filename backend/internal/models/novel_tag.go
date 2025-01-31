package models

type NovelTag struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	TagID   *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"tagId"`
}