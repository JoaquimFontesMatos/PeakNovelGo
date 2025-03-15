package models

// Tag represents a tag in the system. It's used for categorizing content.
//
// Fields:
//   - ID (uint): The unique identifier for the tag.
//   - Name (string): The name of the tag (maximum 255 characters, unique).
//   - Description (string): A description of the tag (maximum 1000 characters).
type Tag struct {
	ID uint `gorm:"primarykey" json:"id"`
	Name string `gorm:"size:255;uniqueIndex;not null" json:"name"`
	Description string `gorm:"size:1000;not null" json:"description"`
}
