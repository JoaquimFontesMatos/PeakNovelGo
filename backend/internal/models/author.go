package models

// Author represents an author.
//
// Fields:
//   - ID (uint): The unique identifier for the author.
//   - Name (string): The name of the author (max 255 characters, unique, and not nullable).
type Author struct {
	ID uint `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:255;uniqueIndex;not null" json:"name"`
}
