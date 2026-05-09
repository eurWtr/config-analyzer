package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		filename    string
		expectError bool
	}{
		{
			name:        "Valid JSON with extension",
			input:       `{"server": {"port": 8080}}`,
			filename:    "config.json",
			expectError: false,
		},
		{
			name:        "Valid YAML without extension",
			input:       "server:\n  port: 8080\n",
			filename:    "",
			expectError: false,
		},
		{
			name:        "Invalid syntax",
			input:       `{"server": {"port": 8080}`, // missing closing brace
			filename:    "config.json",
			expectError: true,
		},
		{
			name:        "JSON stream detection",
			input:       `   {"key": "value"}`, // ensure leading spaces are handled
			filename:    "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := strings.NewReader(tt.input)
			result, err := Parse(reader, tt.filename)

			if tt.expectError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("did not expect error, got: %v", err)
			}

			// If no error, check that the data was parsed
			if !tt.expectError && result == nil {
				t.Errorf("expected result map, got nil")
			}
		})
	}
}
