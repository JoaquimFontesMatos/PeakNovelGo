package dtos

type TTSRequest struct {
    Text      string `json:"text"`
    Rate      int    `json:"rate"`
    Voice     string `json:"voice"`
    NovelID   uint   `json:"novelId"`
    ChapterNo uint   `json:"chapterNo"`
}
