package controllers

import (
	"backend/internal/models"
	"backend/internal/services/interfaces"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// logHandler handles POST requests and writes log entries to a file.
type LogController struct {
	logFilePath string
	logService  interfaces.LogServiceInterface
}

func NewLogController(filePath string, logService interfaces.LogServiceInterface) *LogController {
	return &LogController{
		logFilePath: filePath,
		logService:  logService,
	}
}

func (lc *LogController) SaveLog(c *gin.Context) {
	var entryDTO models.ImportedLogEntry
	// Bind JSON from request to LogEntry struct.
	if err := c.ShouldBindJSON(&entryDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// If no timestamp is provided, set it to current time.
	if entryDTO.Timestamp == "" {
		entryDTO.Timestamp = time.Now().Format(time.RFC3339)
	}

	// Marshal the entry back to JSON for writing.
	logLine, err := json.Marshal(entryDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal log entry dto"})
		return
	}
	logLine = append(logLine, '\n')

	// Open the log file in append mode (create if it doesn't exist).
	f, err := os.OpenFile(lc.logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open log file"})
		return
	}
	defer f.Close()

	// Write the log entry to the file.
	if _, err = f.Write(logLine); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write log entry"})
		return
	}

	logEntry, err := entryDTO.ConvertToModel()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert log entry dto to model"})
		return
	}

	err = lc.logService.CreateLogEntry(logEntry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with success.
	c.JSON(http.StatusOK, gin.H{"message": "Log was created"})
}

// GetLogs handles GET /logs
func (lc *LogController) GetLogs(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	logs, total, err := lc.logService.GetLogs(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Convert each LogEntry model to an ImportedLogEntry DTO
	dtoLogs := make([]models.ImportedLogEntry, 0, len(logs))
	for _, log := range logs {
		dto, convErr := log.ConvertToDTO()
		if convErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": convErr.Error()})
			return
		}
		dtoLogs = append(dtoLogs, dto)
	}

	// Build response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":       dtoLogs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}

// GetLogsByLevel handles GET /logs/level/:level
func (lc *LogController) GetLogsByLevel(c *gin.Context) {
	level := c.Param("level")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) // Default to page 1
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) // Default to 10 items per page
	if err != nil || limit < 1 {
		limit = 10
	}

	logs, total, err := lc.logService.GetLogsByLevel(level, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	if int64(page) > totalPages {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page out of range"})
		return
	}

	// Convert each LogEntry model to an ImportedLogEntry DTO
	dtoLogs := make([]models.ImportedLogEntry, 0, len(logs))
	for _, log := range logs {
		dto, convErr := log.ConvertToDTO()
		if convErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": convErr.Error()})
			return
		}
		dtoLogs = append(dtoLogs, dto)
	}

	// Build response with pagination metadata
	c.JSON(http.StatusOK, gin.H{
		"data":       dtoLogs,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
	})
}
