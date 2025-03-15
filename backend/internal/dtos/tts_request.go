package dtos

// TTSRequest represents a request to generate text-to-speech audio.
//
// Fields:
//   - Text (string): The text to be converted to speech.
//   - Rate (int): The speech rate.
//   - Voice (string): The voice to be used.
//   - NovelID (uint): The ID of the novel.
//   - ChapterNo (uint): The chapter number.
type TTSRequest struct {
	Text      string `json:"text"`
	Rate      int    `json:"rate"`
	Voice     string `json:"voice"`
	NovelID   uint   `json:"novelId"`
	ChapterNo uint   `json:"chapterNo"`
}
