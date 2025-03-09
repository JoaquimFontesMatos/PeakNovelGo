package models

// NovelAuthor represents a many-to-many relationship between novels and authors.
//
// Fields:
//   - NovelID (*uint): The ID of the novel.  Can be NULL.
//   - AuthorID (*uint): The ID of the author. Can be NULL.
type NovelAuthor struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	AuthorID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"authorId"`
}
