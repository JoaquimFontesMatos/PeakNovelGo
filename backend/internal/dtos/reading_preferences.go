package dtos

type ReadingPreferences struct {
    AtomicReading bool           `json:"atomicReading,omitempty"`
    Font          string         `json:"font,omitempty"`
    Theme         string         `json:"theme,omitempty"`
    TTS           TTSPreferences `json:"tts,omitempty"`
}

type TTSPreferences struct {
    Autoplay bool   `json:"autoplay,omitempty"`
    Voice    string `json:"voice,omitempty"`
    Rate     int    `json:"rate,omitempty"`
}
