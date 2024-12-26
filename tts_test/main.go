package main

import (
	"log"

	"github.com/surfaceyu/edge-tts-go/edgeTTS"
)

func main() {
	args := edgeTTS.Args{
		WriteMedia:     "C:/Users/joaquim/Desktop/PeakNovelGo/backend/output.mp3",
	}
	// Initialize the Edge TTS client
	tts := edgeTTS.NewTTS(args)

	edgeTTS.PrintVoices("en-US")

	tts.AddTextWithVoice("Hello! This is a test using Edge TTS in Go.", "en-US-AriaNeural")

	tts.AddTextWithVoice("Hello! This is a test using Edge TTS in Go.", "en-US-JennyNeural")

	tts.AddTextWithVoice("Hello! This is a test using Edge TTS in Go.", "en-US-GuyNeural")

	if tts == nil {
		log.Fatalf("Failed to create TTS client")
	}

	log.Println("Starting speech synthesis...")
	tts.Speak()
	log.Println("Speech synthesis completed.")
}
