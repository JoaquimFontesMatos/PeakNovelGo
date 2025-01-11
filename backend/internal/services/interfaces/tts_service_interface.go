package interfaces

type TTSServiceInterface interface {
	GenerateTTSMap(text, voice, baseURL string) (map[string]string, error)
}
