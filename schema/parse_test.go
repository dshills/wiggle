package schema

import (
	"reflect"
	"testing"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

func TestParseJSONSchema(t *testing.T) {
	// Define a sample JSON schema
	jsonSchema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"field": map[string]interface{}{
				"type":      "string",
				"minLength": 4.0, // Note that in JSON, numbers are float64
				"maxLength": 10.0,
			},
			"age": map[string]interface{}{
				"type":     "integer",
				"minValue": 0.0,
			},
		},
		"required": []interface{}{"field", "age"},
	}

	// Expected schema
	expectedSchema := Schema{
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

	// Call the function to parse the schema
	parsedSchema, err := ParseJSONSchema(jsonSchema)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	// Check that the parsed schema matches the expected schema
	if !reflect.DeepEqual(parsedSchema, expectedSchema) {
		t.Errorf("Parsed schema does not match expected schema. Got %+v, expected %+v", parsedSchema, expectedSchema)
	}
}
