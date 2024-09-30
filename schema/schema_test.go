package schema

import (
	"encoding/json"
	"reflect"
	"testing"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

func TestSchemaToJSON(t *testing.T) {
	// Define a sample Schema struct
	schema := Schema{
		Type: "object",
		Properties: map[string]Schema{
			"field": {
				Type:      "string",
				MinLength: intPointer(4),
				MaxLength: intPointer(10),
			},
			"age": {
				Type:     "integer",
				MinValue: intPointer(0),
			},
		},
		Required: []string{"field", "age"},
	}

	// Expected JSON output
	expectedJSON := `{
		"type": "object",
		"required": [
		  "field",
		  "age"
		],
		"properties": {
		  "field": {
			"type": "string",
			"minLength": 4,
			"maxLength": 10
		  },
		  "age": {
			"type": "integer",
			"minValue": 0
		  }
		}
	  }`

	// Convert the expected JSON string into a map for comparison
	var expectedMap map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedMap); err != nil {
		t.Fatalf("Error unmarshalling expected JSON: %v", err)
	}

	// Convert the Schema struct to JSON
	jsonString, err := schema.ToJSON()
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Convert the generated JSON string into a map for comparison
	var generatedMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &generatedMap); err != nil {
		t.Fatalf("Error unmarshalling generated JSON: %v", err)
	}

	// Compare the two maps
	if !reflect.DeepEqual(generatedMap, expectedMap) {
		t.Errorf("Expected JSON structure:\n%v\nGot:\n%v", expectedMap, generatedMap)
	}
}
