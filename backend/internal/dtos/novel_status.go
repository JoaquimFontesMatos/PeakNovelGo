package dtos

// NovelStatus represents the status of a chapter.
//
// Fields:
//   - NovelUpdatesId (string): The novel title id.
//   - Status (string): The status of the novel (e.g., "completed", "downloading").
type NovelStatus struct {
	NovelUpdatesId string `json:"novelUpdatesId"`
	Status         string `json:"status"`
}
