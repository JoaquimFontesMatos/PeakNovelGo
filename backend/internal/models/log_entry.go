package models

import "encoding/json"

// LogEntry defines the structure of a log entry.
type LogEntry struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Context   string `json:"context,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ImportedLogEntry defines the structure of a log entry.
type ImportedLogEntry struct {
	ID        uint                   `json:"id,omitempty"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Error     json.RawMessage        `json:"error,omitempty"`
}

// ConvertToDTO converts a LogEntry model to an ImportedLogEntry DTO.
func (log *LogEntry) ConvertToDTO() (ImportedLogEntry, error) {
	var contextMap map[string]interface{}
	if log.Context != "" {
		if err := json.Unmarshal([]byte(log.Context), &contextMap); err != nil {
			return ImportedLogEntry{}, err
		}
	}

	var rawError json.RawMessage
	if log.Error != "" {
		if err := json.Unmarshal([]byte(log.Error), &rawError); err != nil {
			return ImportedLogEntry{}, err
		}
	}

	return ImportedLogEntry{
		ID:        log.ID,
		Level:     log.Level,
		Message:   log.Message,
		Timestamp: log.Timestamp,
		Context:   contextMap,
		Error:     rawError,
	}, nil
}

// ConvertToModel converts an ImportedLogEntry DTO to a LogEntry model.
func (dto *ImportedLogEntry) ConvertToModel() (LogEntry, error) {
	// Marshal Context field (map) to a JSON string
	contextBytes, err := json.Marshal(dto.Context)
	if err != nil {
		return LogEntry{}, err
	}

	// Marshal Error field (json.RawMessage) to a JSON string
	errorBytes, err := json.Marshal(dto.Error)
	if err != nil {
		return LogEntry{}, err
	}

	return LogEntry{
		ID:        dto.ID,
		Level:     dto.Level,
		Message:   dto.Message,
		Timestamp: dto.Timestamp,
		Context:   string(contextBytes),
		Error:     string(errorBytes),
	}, nil
}
