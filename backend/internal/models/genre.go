package models

type Genre struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:1000;not null" json:"description"`
}
