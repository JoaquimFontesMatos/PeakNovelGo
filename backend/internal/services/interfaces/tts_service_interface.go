package interfaces

import (
    "backend/internal/dtos"
    "backend/internal/services"

    "github.com/lib-x/edgetts"
)

type TTSServiceInterface interface {
    GenerateTTSMap(ttsRequest *dtos.TTSRequest, baseURL string) ([]services.Paragraph, error)
    GetVoices() ([]edgetts.Voice, error)
}
