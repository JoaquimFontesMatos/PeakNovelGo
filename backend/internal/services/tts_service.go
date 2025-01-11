package services

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/surfaceyu/edge-tts-go/edgeTTS"
)

type TTSService struct {
	OutputDir string
}

// GenerateTTSFile generates TTS audio for a given text and voice and saves it to a file.
func (s *TTSService) GenerateTTSFile(text, voice string) (string, error) {
	// Generate a unique filename based on the current time
	filename := fmt.Sprintf("%d.wav", time.Now().UnixNano())
	filePath := filepath.Join(s.OutputDir, filename)

	args := edgeTTS.Args{
		WriteMedia: filePath,
	}

	tts := edgeTTS.NewTTS(args)
	if tts == nil {
		return "", fmt.Errorf("failed to create TTS client")
	}

	tts.AddText(text, voice, "", "")

	tts.Speak()

	return filename, nil
}

// GenerateParagraphs splits the text into paragraphs based on double newlines.
func (s *TTSService) GenerateParagraphs(text string) []string {
	paragraphs := strings.Split(text, "\n\n")
	for i, paragraph := range paragraphs {
		paragraphs[i] = strings.TrimSpace(paragraph)
	}
	return paragraphs
}

// GenerateTTSMap generates a map of paragraphs to their respective TTS file URLs in parallel, preserving order.
func (s *TTSService) GenerateTTSMap(text, voice, baseURL string) (map[string]string, error) {
	paragraphs := s.GenerateParagraphs(text)

	type result struct {
		Index     int
		Paragraph string
		URL       string
		Err       error
	}

	results := make([]result, len(paragraphs))
	var wg sync.WaitGroup
	var mu sync.Mutex // To protect shared access to results

	for i, paragraph := range paragraphs {
		if paragraph == "" {
			continue
		}

		wg.Add(1)
		go func(i int, paragraph string) {
			defer wg.Done()

			filename, err := s.GenerateTTSFile(paragraph, voice)
			url := ""
			if err == nil {
				url = fmt.Sprintf("%s/%s", baseURL, filename)
			}

			// Lock results slice update
			mu.Lock()
			results[i] = result{
				Index:     i,
				Paragraph: paragraph,
				URL:       url,
				Err:       err,
			}
			mu.Unlock()
		}(i, paragraph)
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Construct the final map, preserving order
	ttsMap := make(map[string]string)
	for _, res := range results {
		if res.Err == nil && res.Paragraph != "" {
			ttsMap[res.Paragraph] = res.URL
		}
	}

	return ttsMap, nil
}
