package rules

import (
	"config-analyzer/internal/models"
	"testing"
)

func TestBindAddressRule_Check(t *testing.T) {
	rule := &BindAddressRule{}

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectedCount int
	}{
		{
			name: "Safe localhost bind",
			config: map[string]interface{}{
				"server": map[string]interface{}{
					"bind": "127.0.0.1",
				},
			},
			expectedCount: 0,
		},
		{
			name: "Unsafe 0.0.0.0 bind",
			config: map[string]interface{}{
				"server": map[string]interface{}{
					"listen_address": "0.0.0.0",
				},
			},
			expectedCount: 1,
		},
		{
			name: "No bind address provided",
			config: map[string]interface{}{
				"app": map[string]interface{}{
					"name": "my-app",
				},
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

			if tt.expectedCount > 0 && issues[0].Severity != models.MEDIUM {
				t.Errorf("expected Severity MEDIUM, got %s", issues[0].Severity)
			}
		})
	}
}
