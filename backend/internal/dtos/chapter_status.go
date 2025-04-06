package dtos

// ChapterStatus represents the status of a chapter.
//
// Fields:
//   - ChapterNo (int): The chapter number.
//   - Status (string): The status of the chapter (e.g., "completed", "downloading").
type ChapterStatus struct {
	ChapterNo int    `json:"chapterNo"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}
