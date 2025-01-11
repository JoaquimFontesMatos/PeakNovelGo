package controllers

import (
	"backend/internal/services/interfaces"
	"backend/internal/validators"
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
	var request struct {
		Text  string `json:"text"`
		Voice string `json:"voice"`
	}

	if err := validators.ValidateBody(ctx, &request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	ttsMap, err := t.ttsService.GenerateTTSMap(request.Text, request.Voice, "http://localhost:8080/tts-files")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate TTS"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"ttsMap": ttsMap})
}
