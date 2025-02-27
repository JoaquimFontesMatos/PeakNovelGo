package controllers

import (
	"backend/internal/dtos"
	"backend/internal/services/interfaces"
	"backend/internal/validators"
	"fmt"
	"log"
	"net/http"
	"os"

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
		log.Println(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	backendDir := fmt.Sprintf("%s/tts-files", os.Getenv("BACKEND_URL"))

	paragraphs, err := t.ttsService.GenerateTTSMap(&request, backendDir)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TTS"})
		return
	}

	ctx.JSON(http.StatusOK, paragraphs)
}

func (t *TTSController) GetVoices(ctx *gin.Context) {
	voices, err := t.ttsService.GetVoices()
	if err != nil {
		log.Println(err.Error())

		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get voices"})
		return
	}

	ctx.JSON(http.StatusOK, voices)
}
