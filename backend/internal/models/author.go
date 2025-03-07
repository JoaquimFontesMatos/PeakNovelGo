package models

type Author struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:255;uniqueIndex;not null" json:"name"`
}