package services

import (
	"backend/internal/dtos"
	"backend/internal/types"
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

// Define supported voices as a map for fast lookup
var supportedVoices = map[string]bool{
	"en-US-AriaNeural":        true,
	"en-GB-LibbyNeural":       true,
	"es-ES-ElviraNeural":      true,
	"en-US-SteffanNeural":     true,
	"en-US-JennyNeural":       true,
	"en-US-MichelleNeural":    true,
	"en-US-EricNeural":        true,
	"en-US-ChristopherNeural": true,
	"en-US-AnaNeural":         true,
	"en-GB-SoniaNeural":       true,
	"en-US-AvaNeural":         true,

	// Add more voices as needed
}

// isVoiceSupported checks if the given voice is supported.
func isVoiceSupported(voice string) bool {
	_, exists := supportedVoices[voice]
	return exists
}

// GenerateTTSFile generates TTS audio for a given text and voice and saves it to a file.
func (s *TTSService) GenerateTTSFile(paragraphs []Paragraph, voice string, rate int) error {
	// Validate the rate value
	if rate < -100 || rate > 100 {
		return fmt.Errorf("rate must be between -100 and 100 (inclusive)")
	}

	// Check if the voice is supported
	if !isVoiceSupported(voice) {
		return fmt.Errorf("invalid voice: %s", voice)
	}

	// Format the rate as a string with a % sign
	rateStr := fmt.Sprintf("%+d%%", rate)

	// Generate a unique filename based on the current time
	options := []edgetts.Option{
		edgetts.WithVoice(voice),
		edgetts.WithRate(rateStr),
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

	go s.scheduleCleanup(paragraphs, 60*time.Minute, 30*time.Minute)

	return nil
}

func (s *TTSService) generateFile(filePath string) (io.Writer, error) {
	// Check if file already exists before creating it
	if _, err := os.Stat(filePath); err == nil {
		return nil, fmt.Errorf("file %s already exists", filePath)
	}
	return os.Create(filePath)
}

func (s *TTSService) GenerateParagraphs(ttsRequest *dtos.TTSRequest, baseUrl string) []Paragraph {
	// Split the text into paragraphs based on double newlines
	paragraphs := strings.Split(ttsRequest.Text, "\n")

	// Process paragraphs to handle dots or other special cases
	processedParagraphs := make([]string, 0, len(paragraphs))
	for _, paragraph := range paragraphs {
		trimmedParagraph := strings.TrimSpace(paragraph)
		trimmedParagraph = strings.ReplaceAll(trimmedParagraph, "<", "")
		trimmedParagraph = strings.ReplaceAll(trimmedParagraph, ">", "")
		if trimmedParagraph == "" || trimmedParagraph == "\n" {
			// Skip empty paragraphs
			continue
		} else if isOnlyDots(trimmedParagraph) {
			// Replace dots with a meaningful placeholder for TTS
			processedParagraphs = append(processedParagraphs, ". . .    scene change    . . .")
		} else {
			// Keep the paragraph as is
			processedParagraphs = append(processedParagraphs, trimmedParagraph)
		}
	}

	// Create the result array
	result := make([]Paragraph, len(processedParagraphs))
	for i, paragraph := range processedParagraphs {
		name := fmt.Sprintf("%d", i)

		filename := fmt.Sprintf("novel_%d_chap%d_%s_[%s_%d].wav", ttsRequest.NovelID, ttsRequest.ChapterNo, name, ttsRequest.Voice, ttsRequest.Rate)
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
		if char != '.' && char != '…' && char != '*' && char != ' ' {
			return false
		}
	}
	return true
}

// GenerateTTSMap generates a map of paragraphs to their respective TTS file URLs in parallel, preserving order.
func (s *TTSService) GenerateTTSMap(ttsRequest *dtos.TTSRequest, baseURL string) ([]Paragraph, error) {
	err := os.MkdirAll(s.OutputDir, 0755)
	if err != nil {
		return nil, types.WrapError(types.INTERNAL_SERVER_ERROR, "An error occurred while creating tts directory", err)
	}

	paragraphs := s.GenerateParagraphs(ttsRequest, baseURL)

	err = s.GenerateTTSFile(paragraphs, ttsRequest.Voice, ttsRequest.Rate)

	if err != nil {
		return nil, err
	}

	return paragraphs, nil
}

func (s *TTSService) GetVoices() ([]edgetts.Voice, error) {
	// Get all voices
	voices, err := edgetts.NewVoiceManager().ListVoices()
	if err != nil {
		return nil, err
	}

	// Filter voices to include only English locales
	var englishVoices []edgetts.Voice
	for _, voice := range voices {
		if isEnglishVoice(voice.Locale) {
			englishVoices = append(englishVoices, voice)
		}
	}

	return englishVoices, nil
}

// isEnglishVoice checks if the locale is for an English voice.
func isEnglishVoice(locale string) bool {
	// English locales start with "en-"
	return len(locale) >= 3 && locale[:3] == "en-"
}

// scheduleCleanup removes only the specific TTS files after a specified duration.
func (s *TTSService) scheduleCleanup(paragraphs []Paragraph, initialDuration time.Duration, extension time.Duration) {
	go func() {
		deadline := time.Now().Add(initialDuration) // Set initial deadline

		for {
			time.Sleep(10 * time.Minute) // Check every 10 minutes

			extend := false
			for _, paragraph := range paragraphs {
				fileInfo, err := os.Stat(paragraph.Filepath)
				if err != nil {
					if os.IsNotExist(err) {
						continue // File already deleted, skip it
					}
					log.Printf("Error accessing file %s: %v", paragraph.Filepath, err)
					continue
				}

				// Get last access time (platform-dependent)
				atime := fileInfo.ModTime() // Linux/macOS: may need syscall to get exact atime

				// Extend if the file was accessed after the deadline was set
				if atime.After(deadline) {
					deadline = time.Now().Add(extension)
					extend = true
					log.Printf("Extended cleanup for %s, new deadline: %v", paragraph.Filepath, deadline)
				}
			}

			if !extend && time.Now().After(deadline) {
				break // Exit loop if no files were accessed and time has passed
			}
		}

		// Delete files after final deadline
		for _, paragraph := range paragraphs {
			err := os.Remove(paragraph.Filepath)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("Failed to delete file %s: %v", paragraph.Filepath, err)
			} else {
				log.Printf("Deleted TTS file: %s", paragraph.Filepath)
			}
		}
	}()
}
