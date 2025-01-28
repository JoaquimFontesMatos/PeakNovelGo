package interfaces

import (
    "backend/internal/services"

    "github.com/lib-x/edgetts"
)

type TTSServiceInterface interface {
    GenerateTTSMap(text string, voice string, novelID uint, chapterNo uint, baseURL string) ([]services.Paragraph, error)
    GetVoices() ([]edgetts.Voice, error)
}
