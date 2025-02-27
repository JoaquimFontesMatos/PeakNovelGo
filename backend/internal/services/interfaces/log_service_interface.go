package interfaces

import (
	"backend/internal/models"
)

type LogServiceInterface interface {
	CreateLogEntry(entry models.LogEntry) error
	GetLogs(page, limit int) ([]models.LogEntry, int64, error)
	GetLogsByLevel(level string, page, limit int) ([]models.LogEntry, int64, error)
}
