package schema

/*
Written by GPT-4o
Directed, modified, and tested by Davin Hills
*/

// ParseJSONSchema takes a map[string]interface{} representation of a JSON schema
// and recursively converts it into a Schema struct. The function handles common JSON schema
// fields such as "type", "properties", "enum", "oneOf", "anyOf", "allOf", "not", as well
// as validation constraints like "minLength", "maxLength", "minValue", and "pattern".
//
// Fields like "properties" are recursively parsed to handle nested objects, while "oneOf",
// "anyOf", and "allOf" allow for combining multiple schemas.
//
// Example Usage:
//
//	jsonSchema := map[string]interface{}{
//	    "type": "object",
//	    "properties": map[string]interface{}{
//	        "field": map[string]interface{}{
//	            "type": "string",
//	            "minLength": 4,
//	            "maxLength": 10,
//	        },
//	        "age": map[string]interface{}{
//	            "type": "integer",
//	            "minValue": 0,
//	        },
//	    },
//	    "required": []interface{}{"field", "age"},
//	}
//
//	schema, err := ParseJSONSchema(jsonSchema)
//	if err != nil {
//	    log.Fatalf("Error parsing schema: %v", err)
//	}
//
//	fmt.Printf("Parsed schema: %+v\n", schema)
//
// Returns a populated Schema struct or an error if the input cannot be parsed into a valid schema.
func ParseJSONSchema(data map[string]interface{}) (Schema, error) {
	schema := Schema{}

	// Handle "type"
	if t, ok := data["type"].(string); ok {
		schema.Type = t
	}

	// Handle "pattern"
	if pattern, ok := data["pattern"].(string); ok {
		schema.Pattern = pattern
	}

	// Handle "required"
	if required, ok := data["required"].([]interface{}); ok {
		for _, r := range required {
			if field, ok := r.(string); ok {
				schema.Required = append(schema.Required, field)
			}
		}
	}

	// Handle "properties"
	if properties, ok := data["properties"].(map[string]interface{}); ok {
		schema.Properties = make(map[string]Schema)
		for key, prop := range properties {
			if propMap, ok := prop.(map[string]interface{}); ok {
				subSchema, err := ParseJSONSchema(propMap)
				if err != nil {
					return Schema{}, err
				}
				schema.Properties[key] = subSchema
			}
		}
	}

	// Handle "enum"
	if enumValues, ok := data["enum"].([]interface{}); ok {
		schema.Enum = enumValues
	}

	// Handle "oneOf"
	if oneOf, ok := data["oneOf"].([]interface{}); ok {
		for _, o := range oneOf {
			if oSchema, ok := o.(map[string]interface{}); ok {
				subSchema, err := ParseJSONSchema(oSchema)
				if err != nil {
					return Schema{}, err
				}
				schema.OneOf = append(schema.OneOf, subSchema)
			}
		}
	}

	// Handle "anyOf"
	if anyOf, ok := data["anyOf"].([]interface{}); ok {
		for _, a := range anyOf {
			if aSchema, ok := a.(map[string]interface{}); ok {
				subSchema, err := ParseJSONSchema(aSchema)
				if err != nil {
					return Schema{}, err
				}
				schema.AnyOf = append(schema.AnyOf, subSchema)
			}
		}
	}

	// Handle "allOf"
	if allOf, ok := data["allOf"].([]interface{}); ok {
		for _, a := range allOf {
			if aSchema, ok := a.(map[string]interface{}); ok {
				subSchema, err := ParseJSONSchema(aSchema)
				if err != nil {
					return Schema{}, err
				}
				schema.AllOf = append(schema.AllOf, subSchema)
			}
		}
	}

	// Handle "not"
	if not, ok := data["not"].(map[string]interface{}); ok {
		subSchema, err := ParseJSONSchema(not)
		if err != nil {
			return Schema{}, err
		}
		schema.Not = &subSchema
	}

	// Handle "minValue"
	if minValue, ok := data["minValue"].(float64); ok {
		val := int(minValue) // Convert float64 to int
		schema.MinValue = &val
	}

	// Handle "minLength"
	if minLength, ok := data["minLength"].(float64); ok {
		val := int(minLength)
		schema.MinLength = &val
	}

	// Handle "maxLength"
	if maxLength, ok := data["maxLength"].(float64); ok {
		val := int(maxLength)
		schema.MaxLength = &val
	}

	// Handle "minItems" and "maxItems"
	if minItems, ok := data["minItems"].(float64); ok {
		val := int(minItems)
		schema.MinItems = &val
	}
	if maxItems, ok := data["maxItems"].(float64); ok {
		val := int(maxItems)
		schema.MaxItems = &val
	}

	// Handle "items" (for arrays)
	if items, ok := data["items"].(map[string]interface{}); ok {
		subSchema, err := ParseJSONSchema(items)
		if err != nil {
			return Schema{}, err
		}
		schema.Items = &subSchema
	}

	// Handle "format"
	if format, ok := data["format"].(string); ok {
		schema.Format = format
	}

	return schema, nil
}
