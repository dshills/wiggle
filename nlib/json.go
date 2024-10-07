package nlib

import (
	"encoding/json"
	"fmt"
	"strings"
)

// JSONAttemptUnmarshal cleans and attempts to unmarshal a possibly malformed JSON string
// into the provided data structure. It extracts the JSON part from the input string,
// cleans it, and then tries to unmarshal the cleaned string into the provided data structure.
//
// Parameters:
//   - str: The input string that potentially contains JSON.
//   - data: A pointer to the variable where the unmarshalled data will be stored. The type of data
//     can be any valid Go data structure (e.g., map, struct, slice).
//
// Returns:
//   - error: An error if cleaning or unmarshalling fails.
func JSONAttemptUnmarshal(str string, data any) error {
	// Clean the string
	cleaned, err := cleanJSONString(str)
	if err != nil {
		return err
	}
	// Try to unmarshal the cleaned string
	_, err = tryUnmarshal(cleaned)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(cleaned), data)
}

// JSONAttemptMap cleans and attempts to unmarshal a possibly malformed JSON string into a map.
// It first tries to extract and clean the JSON part from the string, and then attempts to unmarshal it
// into a map[string]interface{}. If the cleaning or unmarshalling fails, it returns an error.
//
// Parameters:
//   - str: The input string that potentially contains a JSON object.
//
// Returns:
//   - map[string]interface{}: The unmarshalled JSON object as a map.
//   - error: An error if cleaning or unmarshalling fails, or if the result is not a JSON object.
func JSONAttemptMap(str string) (map[string]interface{}, error) {
	// Clean the string
	cleaned, err := cleanJSONString(str)
	if err != nil {
		return nil, err
	}
	// Attempt to unmarshal into a map
	result, err := tryUnmarshal(cleaned)
	if err != nil {
		return nil, err
	}

	// Type assertion: ensure it's a map and not an array
	if dataMap, ok := result.(map[string]interface{}); ok {
		return dataMap, nil
	}
	return nil, fmt.Errorf("expected a JSON object but got something else")
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

	// Step 4: Return cleaned JSON string
	return cleaned, nil
}

// tryUnmarshal attempts to unmarshal the cleaned JSON string into a map or an array
func tryUnmarshal(cleaned string) (interface{}, error) {
	var result interface{}
	err := json.Unmarshal([]byte(cleaned), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return result, nil
}
