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
