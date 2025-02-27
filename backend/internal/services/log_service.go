package services

import (
	"backend/internal/models"
	"backend/internal/repositories/interfaces"
)

// LogService struct represents a service for log management.
type LogService struct {
	logRepository interfaces.LogRepositoryInterface
}

// NewLogService creates a new LogService instance
//
// Parameters:
//   - logRepository interfaces.LogRepositoryInterface (LogRepository instance)
//
// Returns:
//   - *LogService (pointer to the LogService instance)
func NewLogService(logRepository interfaces.LogRepositoryInterface) *LogService {
	return &LogService{logRepository: logRepository}
}

// CreateLogEntry creates a new log entry in the database.
//
// Parameters:
//   - entry models.LogEntry (LogEntry struct)
//
// Returns:
//   - error (nil if the log entry was created successfully, otherwise an error)
//   - INTERNAL_SERVER_ERROR if the log entry could not be created
func (l *LogService) CreateLogEntry(entry models.LogEntry) error {
	return l.logRepository.CreateLogEntry(entry)
}

// GetLogs gets a list of logs.
//
// Parameters:
//   - page int (page number)
//   - limit int (limit of logs per page)
//
// Returns:
//   - []models.LogEntry (list of LogEntry structs)
//   - int64 (total number of logs)
//   - INTERNAL_SERVER_ERROR if the logs could not be fetched
//   - NO_LOGS_ERROR if the logs could not be fetched
func (l *LogService) GetLogs(page, limit int) ([]models.LogEntry, int64, error) {
	return l.logRepository.GetLogs(page, limit)
}

// GetLogsByLevel gets a list of logs by level.
//
// Parameters:
//   - level string (level of the logs)
//   - page int (page number)
//   - limit int (limit of logs per page)
//
// Returns:
//   - []models.LogEntry (list of LogEntry structs)
//   - int64 (total number of logs)
//   - INTERNAL_SERVER_ERROR if the logs could not be fetched
//   - NO_LOGS_ERROR if the logs could not be fetched
//   - LOG_LEVEL_NOT_FOUND_ERROR if the log level doesn't exist
func (l *LogService) GetLogsByLevel(level string, page, limit int) ([]models.LogEntry, int64, error) {
	return l.logRepository.GetLogsByLevel(level, page, limit)
}
