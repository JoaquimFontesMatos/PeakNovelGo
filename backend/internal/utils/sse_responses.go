package utils

import (
	"encoding/json"
	"log"
)

func GetStatusJSON(statuses map[int]string) string {
	statusJSON, err := json.Marshal(statuses)
	if err != nil {
		log.Printf("Failed to marshal statuses: %v", err)
		return "{}"
	}
	return string(statusJSON)
}
