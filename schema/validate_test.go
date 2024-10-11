package schema

import (
	"fmt"
	"testing"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

// intPointer is a helper function that returns a pointer to an integer.
// This is useful for schema fields that expect a pointer to an int (e.g., MinValue, MinItems, MaxLength).
// It allows you to easily pass literal integer values as pointers.
// Params:
// - i: The integer value to create a pointer for.
// Returns a pointer to the given integer value.
func intPointer(i int) *int {
	return &i
}

func TestValidateSimpleType(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				Type: SchemaTypeString,
			},
		},
		Required: []string{"field"},
	}

	// Valid string data
	data := "test string"
	parsedJSON := map[string]interface{}{"field": data}

	if err := Validate(parsedJSON, schema, nil); err != nil {
		t.Errorf("Expected valid string, got error: %v", err)
	}

	// Invalid string data, assigning a non-string type (integer)
	dataInvalid := 123.0 // Ensure it's float64 for JSON numbers
	parsedJSON = map[string]interface{}{"field": dataInvalid}
	if err := Validate(parsedJSON, schema, nil); err == nil {
		t.Errorf("Expected error for invalid type, but got none")
	}
}

func TestValidateEnum(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"color": {
				Type: SchemaTypeString,
				Enum: []interface{}{"red", "green", "blue"},
			},
		},
		Required: []string{"color"},
	}

	// Valid value
	data := "red"
	if err := Validate(map[string]interface{}{"color": data}, schema, nil); err != nil {
		t.Errorf("Expected value in enum, got error: %v", err)
	}

	// Invalid value
	dataInvalid := "yellow"
	if err := Validate(map[string]interface{}{"color": dataInvalid}, schema, nil); err == nil {
		t.Errorf("Expected error for value not in enum, but got none")
	}
}

func TestValidateOneOf(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				OneOf: []Schema{
					{Type: SchemaTypeString},
					{Type: SchemaTypeInteger},
				},
			},
		},
		Required: []string{"field"},
	}

	// Valid string value

	data := map[string]interface{}{"field": "test"}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid value for oneOf, %v got error: %v", "test", err)
	}

	// Valid integer value
	dataInt := map[string]interface{}{"field": 123.0}
	if err := Validate(dataInt, schema, nil); err != nil {
		t.Errorf("Expected valid integer for oneOf, %v got error: %v", 123, err)
	}

	// Invalid value (should match neither schema)
	dataInvalid := map[string]interface{}{"field": []interface{}{1, 2, 3}}
	if err := Validate(dataInvalid, schema, nil); err == nil {
		t.Errorf("Expected error for invalid value in oneOf, %v but got none", []interface{}{1, 2, 3})
	}
}

func TestValidateAnyOf(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				AnyOf: []Schema{
					{Type: SchemaTypeString},
					{Type: SchemaTypeInteger},
				},
			},
		},
		Required: []string{"field"},
	}

	// Valid string value
	data := map[string]interface{}{"field": "test"}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid value for anyOf, %v got error: %v", "test", err)
	}

	// Valid integer value
	dataInt := map[string]interface{}{"field": 123.0}
	if err := Validate(dataInt, schema, nil); err != nil {
		t.Errorf("Expected valid integer for anyOf, %v got error: %v", 123.0, err)
	}

	// Invalid value (should match neither schema)
	dataInvalid := map[string]interface{}{"field": []interface{}{1, 2, 3}}
	if err := Validate(dataInvalid, schema, nil); err == nil {
		t.Errorf("Expected error for invalid value in anyOf, %v but got none", []interface{}{1, 2, 3})
	}
}

func TestValidateAllOf(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				AllOf: []Schema{
					{Type: SchemaTypeString, MinLength: intPointer(4)},
					{Type: SchemaTypeString, MaxLength: intPointer(10)},
				},
			},
		},
		Required: []string{"field"},
	}

	// Valid value (meets both minLength and maxLength)
	data := map[string]interface{}{"field": "valid"}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid value for allOf, got error: %v", err)
	}

	// Invalid value (too short, should trigger an error)
	dataInvalid := map[string]interface{}{"field": "no"}
	if err := Validate(dataInvalid, schema, nil); err == nil {
		t.Errorf("Expected error for value not meeting allOf criteria, but got none")
	}
}

func TestValidateNot(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				Not: &Schema{Type: SchemaTypeInteger},
			},
		},
		Required: []string{"field"},
	}

	// Valid value (not an integer)
	data := map[string]interface{}{"field": "string"}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid value for not, got error: %v", err)
	}

	// Invalid value (matches the 'not' schema)
	dataInt := map[string]interface{}{"field": 123.0}
	if err := Validate(dataInt, schema, nil); err == nil {
		t.Errorf("Expected error for value matching 'not', but got none")
	}
}

func TestValidatePattern(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				Type:    SchemaTypeString,
				Pattern: "^[a-zA-Z]+$",
			},
		},
		Required: []string{"field"},
	}

	// Valid value (letters only)
	data := map[string]interface{}{"field": "ValidString"}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid pattern, got error: %v", err)
	}

	// Invalid value (contains digits)
	dataInvalid := map[string]interface{}{"field": "Invalid123"}
	if err := Validate(dataInvalid, schema, nil); err == nil {
		t.Errorf("Expected error for pattern mismatch, but got none")
	}
}

func TestCustomValidator(t *testing.T) {
	customValidator := func(value interface{}) error {
		if str, ok := value.(string); ok && len(str) > 5 {
			return nil
		}
		return fmt.Errorf("custom validation failed")
	}

	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				Type: SchemaTypeString,
			},
		},
		Required: []string{"field"},
	}

	// Valid custom validation
	data := map[string]interface{}{"field": "customString"}
	if err := Validate(data, schema, customValidator); err != nil {
		t.Errorf("Expected custom validation to pass, got error: %v", err)
	}

	// Invalid custom validation
	dataInvalid := map[string]interface{}{"field": "fail"}
	if err := Validate(dataInvalid, schema, customValidator); err == nil {
		t.Errorf("Expected custom validation to fail, but got none")
	}
}

func TestValidateInteger(t *testing.T) {
	schema := Schema{
		Type: SchemaTypeObject,
		Properties: map[string]Schema{
			"field": {
				Type: SchemaTypeInteger,
			},
		},
		Required: []string{"field"},
	}

	// Valid integer (parsed as float64 in JSON)
	data := map[string]interface{}{"field": 123.0}
	if err := Validate(data, schema, nil); err != nil {
		t.Errorf("Expected valid integer, got error: %v", err)
	}

	// Invalid non-integer number
	dataInvalid := map[string]interface{}{"field": 123.456}
	if err := Validate(dataInvalid, schema, nil); err == nil {
		t.Errorf("Expected error for non-integer number, but got none")
	}
}
