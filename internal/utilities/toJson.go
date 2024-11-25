package utilities

import (
    "encoding/json"
    "log"
)

// ToJSON converts an interface to a JSON string.
func ToJSON(data interface{}) string {
    jsonData, err := json.Marshal(data)
    if err != nil {
        log.Printf("Error marshalling to JSON: %v", err)
        return ""
    }
    return string(jsonData)
}
