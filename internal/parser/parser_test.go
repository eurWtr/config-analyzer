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
			input:       `{"server": {"port": 8080}`, // пропущена закрывающая скобка
			filename:    "config.json",
			expectError: true,
		},
		{
			name:        "JSON stream detection",
			input:       `   {"key": "value"}`, // проверяем пропуск пробелов в начале
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

			// Если ошибки нет, проверяем, что данные распарсились
			if !tt.expectError && result == nil {
				t.Errorf("expected result map, got nil")
			}
		})
	}
}
