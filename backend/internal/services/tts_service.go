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
	Text     string
	Index    int
	Filepath string `json:"-"`
	URL      string
}

// GenerateTTSFile generates TTS audio for a given text and voice and saves it to a file.
func (s *TTSService) GenerateTTSFile(paragraphs []Paragraph, voice string) error {
	// Generate a unique filename based on the current time

	options := []edgetts.Option{
		edgetts.WithVoice(voice),
	}

	edgetts, err := edgetts.NewSpeech(options...)
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

			err = edgetts.AddSingleTask(paragraph.Text, ioWriter)
			if err != nil {
				log.Printf("Error adding task for paragraph %d: %v", paragraph.Index, err)
			}
		}(paragraph)
	}

	// Wait for all tasks to finish
	wg.Wait()

	err = edgetts.StartTasks()

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

// GenerateParagraphs splits the text into paragraphs based on double newlines.
func (s *TTSService) GenerateParagraphs(text string, novelID uint, baseUrl string) []Paragraph {

	paragraphs := strings.Split(text, "\n\n")

	result := make([]Paragraph, len(paragraphs))
	for i, paragraph := range paragraphs {
		name := fmt.Sprintf("%d", i)

		filename := fmt.Sprintf("novel_%d_%s.wav", novelID, name)
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

// GenerateTTSMap generates a map of paragraphs to their respective TTS file URLs in parallel, preserving order.
func (s *TTSService) GenerateTTSMap(text string, voice string, novelID uint, baseURL string) ([]Paragraph, error) {
	os.MkdirAll(s.OutputDir, 0755)

	paragraphs := s.GenerateParagraphs(text, novelID, baseURL)

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
