package services

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/lib-x/edgetts"
)

type TTSService struct {
	OutputDir string
}

type Paragraph struct {
	Text     string `json:"text"`
	Index    int    `json:"index"`
	Filepath string `json:"-"`
	URL      string `json:"url"`
}

// GenerateTTSFile generates TTS audio for a given text and voice and saves it to a file.
func (s *TTSService) GenerateTTSFile(paragraphs []Paragraph, voice string) error {
	// Generate a unique filename based on the current time

	options := []edgetts.Option{
		edgetts.WithVoice(voice),
	}

	tts, err := edgetts.NewSpeech(options...)
	if err != nil {
		return fmt.Errorf("failed to create TTS client")
	}

	// Use a WaitGroup to wait for all tasks to complete
	var wg sync.WaitGroup
	for _, paragraph := range paragraphs {
		wg.Add(1)
		go func(paragraph Paragraph) {
			defer wg.Done()

			ioWriter, err := s.generateFile(paragraph.Filepath)
			if err != nil {
				log.Printf("Error generating file for paragraph %d: %v", paragraph.Index, err)
				return
			}

			err = tts.AddSingleTask(paragraph.Text, ioWriter)
			if err != nil {
				log.Printf("Error adding task for paragraph %d: %v", paragraph.Index, err)
			}
		}(paragraph)
	}

	// Wait for all tasks to finish
	wg.Wait()

	err = tts.StartTasks()

	if err != nil {
		return fmt.Errorf("failed to start tasks")
	}

	go s.scheduleCleanup(s.OutputDir, 15*time.Minute)

	return nil
}

func (s *TTSService) generateFile(filePath string) (io.Writer, error) {
	// Check if file already exists before creating it
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("file %s already exists", filePath)
	}
	return os.Create(filePath)
}

func (s *TTSService) GenerateParagraphs(text string, novelID uint, chapterNo uint, baseUrl string) []Paragraph {
	// Split the text into paragraphs based on double newlines
	paragraphs := strings.Split(text, "\n\n")

	// Process paragraphs to handle dots or other special cases
	processedParagraphs := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph)
		if trimmedParagraph == "" {
			// Skip empty paragraphs
			continue
		} else if isOnlyDots(trimmedParagraph) {
			// Replace dots with a meaningful placeholder for TTS
			processedParagraphs = append(processedParagraphs, "<break time='3s'/>")
		} else {
			// Keep the paragraph as is
			processedParagraphs = append(processedParagraphs, trimmedParagraph)
		}
	}

	// Create the result array
	result := make([]Paragraph, len(processedParagraphs))
	for i, paragraph := range processedParagraphs {
		name := fmt.Sprintf("%d", i)

		filename := fmt.Sprintf("novel_%d_chap%d_%s.wav", novelID, chapterNo, name)
		filePath := filepath.Join(s.OutputDir, filename)

		result[i] = Paragraph{
			Text:     paragraph,
			Index:    i,
			Filepath: filePath,
			URL:      fmt.Sprintf("%s/%s", baseUrl, filename),
		}
	}

	log.Printf("Generated %d paragraphs", len(result))

	return result
}

// Helper function to check if a string contains only dots
func isOnlyDots(s string) bool {
	for _, char := range s {
		if char != '.' {
			return false
		}
	}
	return true
}

// GenerateTTSMap generates a map of paragraphs to their respective TTS file URLs in parallel, preserving order.
func (s *TTSService) GenerateTTSMap(text string, voice string, novelID uint, chapterNo uint, baseURL string) ([]Paragraph, error) {
	os.MkdirAll(s.OutputDir, 0755)

	paragraphs := s.GenerateParagraphs(text, novelID, chapterNo, baseURL)

	err := s.GenerateTTSFile(paragraphs, voice)

	if err != nil {
		return nil, err
	}

	return paragraphs, nil
}

func (s *TTSService) GetVoices() ([]edgetts.Voice, error) {
	return edgetts.NewVoiceManager().ListVoices()
}

// scheduleCleanup schedules a cleanup task to remove the directory after a specified duration.
func (s *TTSService) scheduleCleanup(dir string, duration time.Duration) {
	timer := time.NewTimer(duration)
	<-timer.C // Wait until timer expires

	err := os.RemoveAll(dir)
	if err != nil {
		log.Printf("Failed to clean up directory: %v\n", err)
	} else {
		log.Printf("Directory %s cleaned up successfully after %v\n", dir, duration)
	}
}
