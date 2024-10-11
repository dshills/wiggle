package schema

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"
)

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

const (
	SchemaTypeObject  = "object"
	SchemaTypeArray   = "array"
	SchemaTypeDate    = "date"
	SchemaTypeString  = "string"
	SchemaTypeInteger = "integer"
)

// Validate checks the provided JSON data against the given schema, validating types, required fields,
// formats, enumerations, logical schema constraints (oneOf, anyOf, allOf, not), and custom validators.
// It uses validateWithPath to keep track of the field path for better error reporting.
func Validate(data map[string]interface{}, schema Schema, customValidator CustomValidator) error {
	// Start at the root of the JSON object
	return validateWithPath(data, schema, "root", customValidator)
}

// validateWithPath performs validation on the provided value against the schema and keeps track of the field path.
// This helps with error reporting for nested objects, arrays, and complex validation scenarios.
func validateWithPath(value interface{}, schema Schema, path string, customValidator CustomValidator) error {
	//First, validate the value against its schema using validateValue
	if schema.Type != "" {
		if err := validateValue(value, schema); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// If the schema specifies the value is an array, call validateArray
	if schema.Type == SchemaTypeArray {
		if arr, ok := value.([]interface{}); ok {
			return validateArray(arr, schema, path, customValidator)
		}
		return fmt.Errorf("%s: expected array, got %T", path, value)
	}

	// Validate required fields (only applicable for objects)
	if schema.Type == SchemaTypeObject {
		if objMap, ok := value.(map[string]interface{}); ok {
			for _, field := range schema.Required {
				if _, exists := objMap[field]; !exists {
					return fmt.Errorf("%s.%s: required field is missing", path, field)
				}
			}
		} else {
			return fmt.Errorf("%s: expected object, got %T", path, value)
		}
	}

	// Validate object properties (only applies to objects)
	if schema.Type == SchemaTypeObject && schema.Properties != nil {
		if objMap, ok := value.(map[string]interface{}); ok {
			for key, propSchema := range schema.Properties {
				newPath := fmt.Sprintf("%s.%s", path, key)
				if val, exists := objMap[key]; exists {
					if err := validateWithPath(val, propSchema, newPath, customValidator); err != nil {
						return err
					}
				} else if len(propSchema.Required) > 0 {
					return fmt.Errorf("%s: required property '%s' is missing", newPath, key)
				}
			}
		} else {
			return fmt.Errorf("%s: expected object, got %T", path, value)
		}
	}

	// Validate array elements
	if schema.Type == SchemaTypeArray && schema.Items != nil {
		if arr, ok := value.([]interface{}); ok {
			for i, item := range arr {
				newPath := fmt.Sprintf("%s[%d]", path, i)
				if err := validateWithPath(item, *schema.Items, newPath, customValidator); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("%s: expected array, got %T", path, value)
		}
	}

	// Handle enum validation
	if len(schema.Enum) > 0 {
		if err := validateEnum(value, schema.Enum); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// Handle oneOf validation
	if len(schema.OneOf) > 0 {
		if err := validateOneOf(value, schema.OneOf); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// Handle anyOf validation
	if len(schema.AnyOf) > 0 {
		if err := validateAnyOf(value, schema.AnyOf); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// Handle allOf validation
	if len(schema.AllOf) > 0 {
		if err := validateAllOf(value, schema.AllOf); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// Handle not validation
	if schema.Not != nil {
		if err := validateNot(value, *schema.Not); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	// Validate string patterns
	if schema.Pattern != "" {
		if strVal, ok := value.(string); ok {
			if err := validatePattern(strVal, schema.Pattern); err != nil {
				return fmt.Errorf("%s: %v", path, err)
			}
		} else {
			return fmt.Errorf("%s: expected string for pattern validation, got %T", path, value)
		}
	}

	// Validate format (e.g., email, date)
	if schema.Format != "" {
		if strVal, ok := value.(string); ok {
			if err := validateFormat(strVal, schema.Format); err != nil {
				return fmt.Errorf("%s: %v", path, err)
			}
		} else {
			return fmt.Errorf("%s: expected string for format validation, got %T", path, value)
		}
	}

	// Call custom validators if provided
	if customValidator != nil {
		if err := validateWithCustom(value, schema, customValidator); err != nil {
			return fmt.Errorf("%s: %v", path, err)
		}
	}

	return nil
}

// validatePattern checks if a string value matches the regular expression pattern defined in the schema.
func validatePattern(value string, pattern string) error {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return fmt.Errorf("invalid pattern: %v", err)
	}
	if !matched {
		return fmt.Errorf("value '%s' does not match pattern '%s'", value, pattern)
	}
	return nil
}

// validateEnum checks if a value is one of the allowed enum values defined in the schema.
func validateEnum(value interface{}, enumValues []interface{}) error {
	for _, enumVal := range enumValues {
		if value == enumVal {
			return nil // Value is in the enum
		}
	}
	return fmt.Errorf("value '%v' is not allowed, expected one of %v", value, enumValues)
}

// validateOneOf ensures the value matches exactly one of the schemas in the oneOf array.
func validateOneOf(value interface{}, schemas []Schema) error {
	matchCount := 0
	for _, schema := range schemas {
		if err := validateWithPath(value, schema, "", nil); err == nil {
			matchCount++
		}
	}
	if matchCount != 1 {
		return fmt.Errorf("value must match exactly one schema in 'oneOf', but matched %d", matchCount)
	}
	return nil
}

// validateAnyOf ensures the value matches at least one of the schemas in the anyOf array.
func validateAnyOf(value interface{}, schemas []Schema) error {
	for _, schema := range schemas {
		if err := validateWithPath(value, schema, "", nil); err == nil {
			return nil // Matches at least one schema
		}
	}
	return fmt.Errorf("value does not match any schema in 'anyOf'")
}

// validateAllOf ensures the value matches all of the schemas in the allOf array.
func validateAllOf(value interface{}, schemas []Schema) error {
	for _, schema := range schemas {
		if err := validateWithPath(value, schema, "", nil); err != nil {
			return fmt.Errorf("value does not match all schemas in 'allOf': %v", err)
		}
	}
	return nil
}

// validateNot checks that the value does not match the schema in the not field.
func validateNot(value interface{}, schema Schema) error {
	if err := validateWithPath(value, schema, "", nil); err == nil {
		return fmt.Errorf("value matches disallowed schema in 'not'")
	}
	return nil
}

type CustomValidator func(value interface{}) error

// validateWithCustom applies a custom validator function to the value for additional validation.
func validateWithCustom(value interface{}, schema Schema, customValidator CustomValidator) error {
	// Ensure custom validation is only called for non-object types
	if schema.Type == SchemaTypeObject {
		return nil // Do not run custom validation on objects
	}

	// Call custom validation on the individual field value
	if customValidator != nil {
		if err := customValidator(value); err != nil {
			return err
		}
	}

	return nil
}

// validateArray validates an array of values based on the provided schema.
// It checks if the array length meets the minimum and maximum length constraints (MinItems, MaxItems).
// Additionally, it validates each item in the array against the specified item schema.
// Params:
// - value: A slice of interface{} representing the array to be validated.
// - schema: The schema that defines constraints for the array and its items.
// Returns an error if the array or any of its items fail validation, or nil if validation is successful.
func validateArray(value []interface{}, schema Schema, path string, customValidator CustomValidator) error {
	// Check minItems and maxItems constraints
	if schema.MinItems != nil && len(value) < *schema.MinItems {
		return fmt.Errorf("%s: array has fewer than %d items", path, *schema.MinItems)
	}
	if schema.MaxItems != nil && len(value) > *schema.MaxItems {
		return fmt.Errorf("%s: array has more than %d items", path, *schema.MaxItems)
	}

	// Validate each item in the array against the `Items` schema
	if schema.Items != nil {
		for i, item := range value {
			newPath := fmt.Sprintf("%s[%d]", path, i)
			if err := validateWithPath(item, *schema.Items, newPath, customValidator); err != nil {
				return fmt.Errorf("%s: %v", newPath, err)
			}
		}
	}
	return nil
}

// validateFormat checks if a string value adheres to a specific format, such as "email" or "date".
// Currently supported formats include:
// - "email": Validates that the string contains an "@" symbol as a basic email format check.
// - "date": Validates that the string matches the "YYYY-MM-DD" format using the Go `time.Parse` function.
// Params:
// - value: A string value that needs to be validated for a specific format.
// - format: A string specifying the format to validate against (e.g., "email", "date").
// Returns an error if the value does not match the expected format, or nil if validation is successful.
func validateFormat(value string, format string) error {
	switch format {
	case "email":
		// Basic email format validation (you could replace this with a more sophisticated regex)
		if !strings.Contains(value, "@") {
			return fmt.Errorf("invalid email format")
		}
	case "date":
		_, err := time.Parse("2006-01-02", value) // Expected format is YYYY-MM-DD
		if err != nil {
			return fmt.Errorf("invalid date format, expected YYYY-MM-DD")
		}
		// Add other formats like "uri", "ipv4", etc.
	}

	return nil
}

// validateValue performs validation on a single value against the provided schema.
// It first checks the type of the value using the `checkType` function. Depending on the type, it applies further validations:
// - For numbers (integer, number), it checks if the value meets the minimum value (MinValue).
// - For strings, it checks if the string length meets the minimum (MinLength) and maximum (MaxLength) length constraints.
// Params:
// - value: The value to be validated (can be of any type: string, number, object, etc.).
// - schema: The schema that defines the expected type and additional constraints (like min/max values).
// Returns an error if the value fails validation, or nil if validation is successful.
func validateValue(value interface{}, schema Schema) error {
	// Check type first
	if err := checkType(value, schema.Type); err != nil {
		return err
	}

	// If it's a number (integer or float), validate MinValue
	if schema.Type == SchemaTypeInteger || schema.Type == "number" {
		if schema.MinValue != nil {
			if number, ok := value.(float64); ok { // JSON numbers are always float64 in Go
				if number < float64(*schema.MinValue) {
					return fmt.Errorf("value %v is less than minimum value of %d", number, *schema.MinValue)
				}
			} else {
				return fmt.Errorf("expected a number, got %v", reflect.TypeOf(value))
			}
		}
	}

	// If it's a string, validate string-specific constraints
	if schema.Type == SchemaTypeString {
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("expected string, got %v", reflect.TypeOf(value))
		}

		if schema.MinLength != nil && len(strValue) < *schema.MinLength {
			return fmt.Errorf("string is shorter than minimum length of %d", *schema.MinLength)
		}

		if schema.MaxLength != nil && len(strValue) > *schema.MaxLength {
			return fmt.Errorf("string is longer than maximum length of %d", *schema.MaxLength)
		}
	}

	return nil
}

// checkType ensures that a value is of the expected type as defined in the schema.
// Supported types include:
// - SchemaTypeString: Ensures the value is a string.
// - SchemaTypeInteger: Ensures the value is a number (JSON numbers are parsed as float64 in Go).
// - SchemaTypeObject: Ensures the value is a map (representing a JSON object).
// - SchemaTypeArray: Ensures the value is a slice (representing a JSON array).
// Params:
// - value: The value whose type needs to be checked.
// - expectedType: A string representing the expected type (e.g., SchemaTypeString, SchemaTypeInteger).
// Returns an error if the value is not of the expected type, or nil if the type check passes.
func checkType(value interface{}, expectedType string) error {
	switch expectedType {
	case SchemaTypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case SchemaTypeInteger:
		// JSON numbers are float64, so we need to treat float64 as valid for integers
		if number, ok := value.(float64); ok {
			// Check if it's a whole number (an integer)
			if number != float64(int(number)) {
				return fmt.Errorf("expected integer, got non-integer number %v", value)
			}
		} else {
			return fmt.Errorf("expected integer, got %T", value)
		}
	case SchemaTypeObject:
		// Objects in JSON are represented as maps in Go (map[string]interface{})
		if _, ok := value.(map[string]interface{}); !ok {
			return fmt.Errorf("expected object, got %T", value)
		}
	case SchemaTypeArray:
		if _, ok := value.([]interface{}); !ok {
			return fmt.Errorf("expected array, got %T", value)
		}
	default:
		return fmt.Errorf("unsupported type '%s'", expectedType)
	}
	return nil
}
