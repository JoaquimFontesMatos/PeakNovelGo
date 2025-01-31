package controllers

import (
	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/validators"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TTSController struct {
	ttsService interfaces.TTSServiceInterface
}

func NewTTSController(ttsService interfaces.TTSServiceInterface) *TTSController {
	return &TTSController{ttsService: ttsService}
}

func (t *TTSController) GenerateTTS(ctx *gin.Context) {
	var request dtos.TTSRequest

	if err := validators.ValidateBody(ctx, &request); err != nil {
		log.Printf(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	paragraphs, err := t.ttsService.GenerateTTSMap(&request, "http://localhost:8081/tts-files")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TTS"})
		return
	}

	ctx.JSON(http.StatusOK, paragraphs)
}

func (t *TTSController) GetVoices(ctx *gin.Context) {
	voices, err := t.ttsService.GetVoices()
	if err != nil {
		log.Printf(err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get voices"})
		return
	}

	ctx.JSON(http.StatusOK, voices)
}
