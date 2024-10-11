package nlib

import (
	"reflect"
	"testing"
)

func TestJSONAttemptMap(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:  "Valid JSON object",
			input: `{"key1": "value1", "key2": 2}`,
			want: map[string]interface{}{
				"key1": "value1",
				"key2": 2.0, // JSON unmarshals numbers as float64
			},
			wantErr: false,
		},
		{
			name:  "Valid JSON object with junk",
			input: `randomjunk{"key1":"value1"}morejunk`,
			want: map[string]interface{}{
				"key1": "value1",
			},
			wantErr: false,
		},
		{
			name:  "Valid JSON object with escaped quotes",
			input: `\"{\"key\":\"value\"}\"`,
			want: map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "Malformed JSON",
			input:   `randomjunk { key: "value" } morejunk`,
			wantErr: true,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "Non-JSON input",
			input:   "this is just text",
			wantErr: true,
		},
		{
			name:    "Valid JSON array",
			input:   `randomtext [ { "key1": "value1" }, { "key2": "value2" } ] more random text`,
			wantErr: true, // We expect an object but get an array
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONAttemptMap(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONAttemptMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSONAttemptMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSONAttemptUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		target  interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name:   "Valid JSON object into map",
			input:  `{"key1": "value1", "key2": 2}`,
			target: &map[string]interface{}{},
			want: &map[string]interface{}{
				"key1": "value1",
				"key2": 2.0, // JSON unmarshals numbers as float64
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON array into slice",
			input:  `[{"key1": "value1"}, {"key2": "value2"}]`,
			target: &[]map[string]interface{}{},
			want: &[]map[string]interface{}{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			wantErr: false,
		},
		{
			name:   "Valid JSON with extra characters",
			input:  `randomtext{"key": "value"}extrajunk`,
			target: &map[string]interface{}{},
			want: &map[string]interface{}{
				"key": "value",
			},
			wantErr: false,
		},
		{
			name:    "Malformed JSON",
			input:   `{"key": "value"`,
			target:  &map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "Empty input",
			input:   "",
			target:  &map[string]interface{}{},
			wantErr: true,
		},
		{
			name:    "Non-JSON input",
			input:   "randomtext",
			target:  &map[string]interface{}{},
			wantErr: true,
		},
		{
			name:   "Valid JSON array with junk",
			input:  `randomtext [ { "key1": "value1" }, { "key2": "value2" } ] morejunk`,
			target: &[]map[string]interface{}{},
			want: &[]map[string]interface{}{
				{"key1": "value1"},
				{"key2": "value2"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := JSONAttemptUnmarshal(tt.input, tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSONAttemptUnmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.target, tt.want) {
				t.Errorf("JSONAttemptUnmarshal() = %v, want %v", tt.target, tt.want)
			}
		})
	}
}

func TestCleanJSONString(t *testing.T) {
	jstr := `
[
  {
    "task_name": "Implement Rope Data Structure",
    "task_number": 1,
    "task_steps": [
      "Define the basic structure of the Rope data type.",
      "Create nodes to represent both leaves (with strings) and composites (concatenation of other nodes).",
      "Incorporate size tracking for each node."
    ],
    "notes": [
      "Focus on memory efficiency and split-forwarding."
    ]
  },
  {
    "task_name": "Implement Insert Method",
    "task_number": 2,
    "task_steps": [
      "Design the function signature to insert a string at a given character index.",
      "Traverse the Rope to find the appropriate insert location.",
      "Balance the Rope if necessary after insertion."
    ],
    "notes": [
      "Think about edge cases such as inserting at the beginning or end."
    ]
  },
  {
    "task_name": "Implement Delete Method",
    "task_number": 3,
    "task_steps": [
      "Design the function signature to remove a substring given a starting index and a length.",
      "Traverse the Rope to locate and remove the specified substring.",
      "Balance the Rope if necessary after deletion."
    ],
    "notes": [
      "Ensure memory optimization during deletion."
    ]
  },
  {
    "task_name": "Implement Substring Method",
    "task_number": 4,
    "task_steps": [
      "Design the function to extract a substring given a starting index and length.",
      "Traverse the Rope efficiently to gather parts of the substring.",
      "Return the extracted substring as a string."
    ],
    "notes": [
      "Optimize traversal for minimal node visitation."
    ]
  },
  {
    "task_name": "Implement Character Index Conversion",
    "task_number": 5,
    "task_steps": [
      "Provide a method to convert a character position to a node-specific index.",
      "Utilize Rope's structure for quick traversal to the target position."
    ],
    "notes": [
      "Consider caching mechanisms to enhance performance."
    ]
  },
  {
    "task_name": "Implement Line and Column Conversion",
    "task_number": 6,
    "task_steps": [
      "Design methods to convert between character index and line-column values.",
      "Implement efficient line tracking during modifications."
    ],
    "notes": [
      "Handle newline characters intelligently."
    ]
  },
  {
    "task_name": "Write Unit Tests for Rope Operations",
    "task_number": 7,
    "task_steps": [
      "Develop test cases for basic operations: insert, delete, substring.",
      "Test additional methods such as conversions between index and line-column."
    ],
    "notes": [
      "Consider edge cases and performance under varying Rope sizes."
    ]
  },
  {
    "task_name": "Optimize Rope Balancing",
    "task_number": 8,
    "task_steps": [
      "Analyze and improve the balance operation after insertion and deletion.",
      "Explore strategies like rebalancing thresholds."
    ],
    "notes": [
      "Aim for consistent performance improvements in common use cases."
    ]
  },
  {
    "task_name": "Document Public API",
    "task_number": 9,
    "task_steps": [
      "Draft comprehensive documentation for each public method.",
      "Include usage examples and expected behavior."
    ],
    "notes": [
      "Emphasize clarity for end-user developers."
    ]
  },
  {
    "task_name": "Benchmark and Profile Rope Implementation",
    "task_number": 10,
    "task_steps": [
      "Create realistic workload scenarios for benchmarking.",
      "Use profiling tools to identify performance bottlenecks."
    ],
    "notes": [
      "Focus on optimizations revealed by profiling data."
    ]
  }
]
`
	jstr = "```go" + jstr + "```"
	str, err := cleanJSONString(jstr)
	if err != nil {
		t.Fatal(err)
	}
	_, err = JSONAttemptArrayMap(str)
	if err != nil {
		t.Fatal(err)
	}
}
