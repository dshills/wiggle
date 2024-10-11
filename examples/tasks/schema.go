package main

var schemaStr = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "tasks": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "task_name": {
            "type": "string",
            "description": "The name of the task"
          },
          "task_number": {
            "type": "integer",
            "description": "The number identifying the task"
          },
          "task_steps": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "The steps to complete the task"
          },
          "notes": {
            "type": "array",
            "items": {
              "type": "string"
            },
            "description": "Additional notes for the task"
          }
        },
        "required": ["task_name", "task_number", "task_steps", "notes"],
        "additionalProperties": false
      }
    }
  },
  "required": ["tasks"],
  "additionalProperties": false
}
`
