package schema

import "encoding/json"

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

// Schema structure definition
type Schema struct {
	Type       string            `json:"type"`
	Pattern    string            `json:"pattern,omitempty"`
	Required   []string          `json:"required,omitempty"`
	Properties map[string]Schema `json:"properties,omitempty"`
	Enum       []interface{}     `json:"enum,omitempty"`
	OneOf      []Schema          `json:"oneOf,omitempty"`
	AnyOf      []Schema          `json:"anyOf,omitempty"`
	AllOf      []Schema          `json:"allOf,omitempty"`
	Not        *Schema           `json:"not,omitempty"`
	MinValue   *int              `json:"minValue,omitempty"`
	MinLength  *int              `json:"minLength,omitempty"`
	MaxLength  *int              `json:"maxLength,omitempty"`
	MinItems   *int              `json:"minItems,omitempty"`
	MaxItems   *int              `json:"maxItems,omitempty"`
	Items      *Schema           `json:"items,omitempty"`
	Format     string            `json:"format,omitempty"`
}

// ToJSON converts a Schema struct to a JSON string
func (schema Schema) ToJSON() (string, error) {
	schemaMap := schema.toMap()
	jsonBytes, err := json.MarshalIndent(schemaMap, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// toMap Converts a Schema struct back into a JSON-compatible map[string]interface{}
func (schema Schema) toMap() map[string]interface{} {
	result := make(map[string]interface{})

	// Set the "type" field if it exists
	if schema.Type != "" {
		result["type"] = schema.Type
	}

	// Set the "pattern" field
	if schema.Pattern != "" {
		result["pattern"] = schema.Pattern
	}

	// Set the "required" field if there are any required fields
	if len(schema.Required) > 0 {
		required := make([]interface{}, len(schema.Required))
		for i, req := range schema.Required {
			required[i] = req
		}
		result["required"] = required
	}

	// Set the "properties" field if there are any properties
	if len(schema.Properties) > 0 {
		properties := make(map[string]interface{})
		for key, prop := range schema.Properties {
			properties[key] = prop.toMap()
		}
		result["properties"] = properties
	}

	// Set the "enum" field if there are enum values
	if len(schema.Enum) > 0 {
		result["enum"] = schema.Enum
	}

	// Set the "oneOf" field if present
	if len(schema.OneOf) > 0 {
		oneOf := make([]interface{}, len(schema.OneOf))
		for i, o := range schema.OneOf {
			oneOf[i] = o.toMap()
		}
		result["oneOf"] = oneOf
	}

	// Set the "anyOf" field if present
	if len(schema.AnyOf) > 0 {
		anyOf := make([]interface{}, len(schema.AnyOf))
		for i, a := range schema.AnyOf {
			anyOf[i] = a.toMap()
		}
		result["anyOf"] = anyOf
	}

	// Set the "allOf" field if present
	if len(schema.AllOf) > 0 {
		allOf := make([]interface{}, len(schema.AllOf))
		for i, a := range schema.AllOf {
			allOf[i] = a.toMap()
		}
		result["allOf"] = allOf
	}

	// Set the "not" field if present
	if schema.Not != nil {
		result["not"] = schema.Not.toMap()
	}

	// Set number fields like "minValue", "minLength", "maxLength"
	if schema.MinValue != nil {
		result["minValue"] = *schema.MinValue
	}
	if schema.MinLength != nil {
		result["minLength"] = *schema.MinLength
	}
	if schema.MaxLength != nil {
		result["maxLength"] = *schema.MaxLength
	}

	// Set array constraints like "minItems", "maxItems"
	if schema.MinItems != nil {
		result["minItems"] = *schema.MinItems
	}
	if schema.MaxItems != nil {
		result["maxItems"] = *schema.MaxItems
	}

	// Set the "items" field if present (for array validation)
	if schema.Items != nil {
		result["items"] = schema.Items.toMap()
	}

	// Set the "format" field if present
	if schema.Format != "" {
		result["format"] = schema.Format
	}

	return result
}
