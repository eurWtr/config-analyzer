package rules

import (
	"config-analyzer/internal/models"
	"testing"
)

func TestDebugModeRule_Check(t *testing.T) {
	rule := &DebugModeRule{}

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectedCount int
	}{
		{
			name: "Boolean debug mode enabled",
			config: map[string]interface{}{
				"app": map[string]interface{}{"debug": true},
			},
			expectedCount: 1,
		},
		{
			name:          "String debug mode enabled",
			config:        map[string]interface{}{"debug": "yes"},
			expectedCount: 1,
		},
		{
			name:          "Debug log level",
			config:        map[string]interface{}{"log_level": "DEBUG"},
			expectedCount: 1,
		},
		{
			name: "Production ready config",
			config: map[string]interface{}{
				"debug":     false,
				"log_level": "info",
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Check(tt.config, "test.yaml")
			if len(issues) != tt.expectedCount {
				t.Errorf("expected %d issues, got %d", tt.expectedCount, len(issues))
			}
			if len(issues) > 0 && issues[0].Severity != models.LOW {
				t.Errorf("expected Severity LOW, got %v", issues[0].Severity)
			}
		})
	}
}
