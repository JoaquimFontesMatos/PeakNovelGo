package utils

import (
	"encoding/json"
	"log"
)

// GetStatusJSON converts a map of integer keys and string values to a JSON string. Returns "{}" if marshaling fails.
//
// Parameters:
// 	- statuses map[int]string (map of the statuses of an SSE response, the integer keys correspond to the ids of the
//	elements in the map, and the string values are the status of those ids)
//
// Returns:
//	- string (map stringified into a json string)
func GetStatusJSON(statuses map[int]string) string {
	statusJSON, err := json.Marshal(statuses)
	if err != nil {
		log.Printf("Failed to marshal statuses: %v", err)
		return "{}"
	}
	return string(statusJSON)
}
