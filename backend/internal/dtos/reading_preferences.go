package dtos

// ReadingPreferences represents the user's preferences for reading.
//
// Fields:
//   - AtomicReading (bool): If true, reading will happen atomically.
//   - Font (string): The preferred font for reading.
//   - Theme (string): The preferred theme for reading.
//   - TTS (TTSPreferences): The preferences for text-to-speech.
type ReadingPreferences struct {
	AtomicReading bool           `json:"atomicReading"`
	Font          string         `json:"font"`
	Theme         string         `json:"theme"`
	TTS           TTSPreferences `json:"tts"`
}

// TTSPreferences represents the user's preferences for text-to-speech.
//
// Fields:
//   - Autoplay (bool): Whether text-to-speech should automatically start playing.
//   - Voice (string): The name or identifier of the desired voice.
//   - Rate (int): The playback rate, typically a percentage of normal speed.
type TTSPreferences struct {
	Autoplay bool   `json:"autoplay"`
	Voice    string `json:"voice"`
	Rate     int    `json:"rate"`
}
