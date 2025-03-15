package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

// GetStatusJSON converts a map of integer keys and string values to a JSON string. Returns "{}" if marshaling fails.
//
// Parameters:
//   - statuses map[int]string (map of the statuses of an SSE response, the integer keys correspond to the ids of the
//     elements in the map, and the string values are the status of those ids)
//
// Returns:
//   - string (map stringified into a json string)
func GetStatusJSON(statuses map[any]string) string {
	converted := make(map[string]string)

	for k, v := range statuses {
		converted[fmt.Sprintf("%v", k)] = v // Convert all keys to string
	}

	statusJSON, err := json.Marshal(converted)
	if err != nil {
		log.Printf("Failed to marshal statuses: %v", err)
		return "{}"
	}
	return string(statusJSON)
}
