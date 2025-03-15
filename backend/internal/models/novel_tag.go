package models

// NovelTag represents a many-to-many relationship between novels and tags.
//
// Fields:
//   - NovelID (*uint): The ID of the novel.  Can be NULL.
//   - TagID (*uint): The ID of the tag. Can be NULL.
type NovelTag struct {
	NovelID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"novelId"`
	TagID *uint `gorm:"constraint:OnDelete:SET NULL;index" json:"tagId"`
}
