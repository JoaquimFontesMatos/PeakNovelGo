package dtos

type ReadingPreferences struct {
    AtomicReading bool           `json:"atomicReading"`
    Font          string         `json:"font"`
    Theme         string         `json:"theme"`
    TTS           TTSPreferences `json:"tts"`
}

type TTSPreferences struct {
    Autoplay bool   `json:"autoplay"`
    Voice    string `json:"voice"`
    Rate     int    `json:"rate"`
}
