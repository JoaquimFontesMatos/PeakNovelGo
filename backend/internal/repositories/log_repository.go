package repositories

import (
	"backend/internal/models"
	"backend/internal/types"
	"log"

	"gorm.io/gorm"
)

// LogRepository struct represents a repository for log management.
type LogRepository struct {
	db *gorm.DB
}

// NewLogRepository creates a new LogRepository instance
//
// Parameters:
//   - db *gorm.DB (Gorm database connection)
//
// Returns:
//   - *LogRepository (pointer to the LogRepository instance)
func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{db: db}
}

// CreateLogEntry creates a new log entry in the database.
//
// Parameters:
//   - entry (log entry)
//
// Returns:
//   - error (nil if the log entry was created successfully, otherwise an error)
//   - INTERNAL_SERVER_ERROR if the log entry could not be created
func (l *LogRepository) CreateLogEntry(entry models.LogEntry) error {
	if err := l.db.Create(&entry).Error; err != nil {
		log.Println(err)
		return types.WrapError("INTERNAL_SERVER_ERROR", "Failed to create log", err)
	}
	return nil
}

func (l *LogRepository) GetLogs(page, limit int) ([]models.LogEntry, int64, error) {
	var logs []models.LogEntry
	var total int64

	// Count total chapters for the novel
	if err := l.db.Model(&models.LogEntry{}).Count(&total).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_LOGS_ERROR, "No logs found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of logs", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := l.db.Model(&models.LogEntry{}).
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_LOGS_ERROR, "No logs found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch logs", err)
	}

	return logs, total, nil
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
func (l *LogRepository) GetLogsByLevel(level string, page, limit int) ([]models.LogEntry, int64, error) {
	var logs []models.LogEntry
	var total int64

	if level != "info" && level != "error" && level != "warning" && level != "debug" {
		return nil, 0, types.WrapError(types.LOG_LEVEL_NOT_FOUND_ERROR, "Log level non existent", nil)
	}

	// Count total logs for the level
	if err := l.db.Model(&models.LogEntry{}).
		Where("level = ?", level).
		Count(&total).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_LOGS_ERROR, "No logs found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to get the total number of logs", err)
	}

	// Apply pagination and ordering
	offset := (page - 1) * limit
	if err := l.db.Model(&models.LogEntry{}).
		Where("level = ?", level).
		Limit(limit).
		Offset(offset).
		Find(&logs).Error; err != nil {

		if err.Error() == "record not found" {
			return nil, 0, types.WrapError(types.NO_LOGS_ERROR, "No logs found", nil)
		}

		return nil, 0, types.WrapError(types.INTERNAL_SERVER_ERROR, "Failed to fetch logs", err)
	}

	return logs, total, nil
}
