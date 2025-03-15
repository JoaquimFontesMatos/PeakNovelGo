package models

import "encoding/json"

// LogEntry represents a log entry.
//
// Fields:
//   - ID (uint): The unique identifier of the log entry.
//   - Level (string): The severity level of the log entry (e.g., "DEBUG", "INFO", "ERROR").
//   - Message (string): The message of the log entry.
//   - Timestamp (string): The timestamp of the log entry.
//   - Context (string): Optional context information related to the log entry.
//   - Error (string): Optional error message associated with the log entry.
type LogEntry struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Context   string `json:"context,omitempty"`
	Error     string `json:"error,omitempty"`
}

// ImportedLogEntry represents a log entry imported from an external source.  It contains the log level, message, timestamp, context, and any associated error.
//
// Fields:
//   - ID (uint): The unique identifier of the log entry.  Omitted when not available.
//   - Level (string): The severity level of the log entry (e.g., "DEBUG", "INFO", "ERROR").
//   - Message (string): The main message of the log entry.
//   - Timestamp (string): The timestamp of when the log entry was created.
//   - Context (map[string]interface{}):  A map containing additional contextual information related to the log entry. Omitted when not available.
//   - Error (json.RawMessage):  The error associated with the log entry, if any.  Stored as a RawMessage to avoid unmarshalling issues with unknown error types. Omitted when no error occurred.
type ImportedLogEntry struct {
	ID        uint                   `json:"id,omitempty"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Context   map[string]interface{} `json:"context,omitempty"`
	Error     json.RawMessage        `json:"error,omitempty"`
}

// ConvertToDTO converts a LogEntry to an ImportedLogEntry.
//
// It handles JSON unmarshalling of the Context and Error fields,
// returning an error if unmarshalling fails.  If the Context or Error fields are empty strings,
// the corresponding fields in the ImportedLogEntry will be nil or empty.
//
// Parameters:
//  - log (*LogEntry): The LogEntry to convert.
//
// Returns:
//   - ImportedLogEntry: The converted ImportedLogEntry.
//   - error: An error if JSON unmarshalling fails.  Nil if successful.
//
// Error types:
//   - json.UnmarshalTypeError: If the JSON in Context or Error fields is invalid.
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
//
// It marshals the Context and Error fields (which are maps and json.RawMessage respectively) into JSON strings before assigning them to the LogEntry struct.  Any marshaling errors are returned.
//
// Parameters:
//  - dto (*ImportedLogEntry): The ImportedLogEntry DTO to convert.
//
// Returns:
//   - LogEntry: The converted LogEntry model.
//   - error: An error encountered during JSON marshaling of the Context or Error fields.  Returns an empty LogEntry if an error occurs.
//
// Error types:
//   - json.MarshalError: If marshaling the Context or Error field fails.
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
