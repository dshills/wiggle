package nlib

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONAttemptUnmarshal cleans and attempts to unmarshal a possibly malformed JSON string
// into the provided data structure. It extracts the JSON part from the input string,
// cleans it, and then tries to unmarshal the cleaned string into the provided data structure.
func JSONAttemptUnmarshal(str string, data any) error {
	// Clean the string
	cleaned, err := cleanJSONString(str)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cleaned), data)
}

// JSONAttemptMap cleans and attempts to unmarshal a possibly malformed JSON string into a map.
// It first tries to extract and clean the JSON part from the string, and then attempts to unmarshal it
// into a map[string]interface{}. If the cleaning or unmarshalling fails, it returns an error.
func JSONAttemptMap(str string) (map[string]interface{}, error) {
	// Clean the string
	cleaned, err := cleanJSONString(str)
	if err != nil {
		return nil, err
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("no json found")
	}
	if cleaned[0] == '[' {
		return nil, fmt.Errorf("not a json object. Maybe an array")
	}

	var mp map[string]interface{}
	err = json.Unmarshal([]byte(cleaned), &mp)
	return mp, err
}

func JSONAttemptArrayMap(str string) ([]map[string]interface{}, error) {
	// Clean the string
	cleaned, err := cleanJSONString(str)
	if err != nil {
		return nil, err
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("no json found")
	}
	if cleaned[0] == '{' {
		return nil, fmt.Errorf("not a json array. Maybe an object")
	}

	var mp []map[string]interface{}
	err = json.Unmarshal([]byte(cleaned), &mp)
	return mp, err
}

func findJSON(input string) string {
	depth := 0
	start := -1
	for i, r := range input {
		switch r {
		case '{', '[':
			if depth == 0 {
				start = i
			}
			depth++
		case '}', ']':
			depth--
			if depth == 0 {
				return input[start : i+1]
			}
		}
	}
	return "" // No valid JSON found
}

// cleanJSONString attempts to clean a possibly malformed JSON string
func cleanJSONString(input string) (string, error) {
	// Step 1: Extract potential JSON part using findJSON to find the curly braces or array brackets
	jsonPart := findJSON(input)
	if jsonPart == "" {
		return "", fmt.Errorf("no valid JSON")
	}

	// Step 2: Handle escaped characters
	// Replace escaped backslashes and quotes (e.g. \" -> ")
	cleaned := strings.ReplaceAll(jsonPart, `\"`, `"`)

	// Step 3: Check for extra backslashes (optional, depends on your input)
	cleaned = strings.ReplaceAll(cleaned, `\\`, `\`)

	// Step 4: Remove any spaces
	cleaned = strings.TrimSpace(cleaned)

	// Step 5: Return cleaned JSON string
	return cleaned, nil
}
