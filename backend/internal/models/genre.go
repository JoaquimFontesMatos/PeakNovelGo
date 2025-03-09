package models

// Genre represents a genre of media.
//
// Fields:
//   - ID (uint): The unique identifier for the genre.
//   - Name (string): The name of the genre. Must be unique.
//   - Description (string): A description of the genre.
type Genre struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:1000;not null" json:"description"`
}
