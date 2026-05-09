package rules

import (
	"config-analyzer/internal/models"
	"testing"
)

func TestPlainPasswordRule_Check(t *testing.T) {
	rule := &PlainPasswordRule{}

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectedCount int
	}{
		{
			name:          "Clear text password",
			config:        map[string]interface{}{"db.password": "supersecret"},
			expectedCount: 1,
		},
		{
			name:          "Vault reference (safe)",
			config:        map[string]interface{}{"api_key": "vault://secret/data"},
			expectedCount: 0,
		},
		{
			name:          "Environment variable (safe)",
			config:        map[string]interface{}{"token": "${GITHUB_TOKEN}"},
			expectedCount: 0,
		},
		{
			name:          "Masked password (safe)",
			config:        map[string]interface{}{"password": "***"},
			expectedCount: 0,
		},
		{
			name:          "Empty password (safe/ignored)",
			config:        map[string]interface{}{"pass": ""},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Check(tt.config, "test.yaml")
			if len(issues) != tt.expectedCount {
				t.Errorf("expected %d issues, got %d", tt.expectedCount, len(issues))
			}
			if len(issues) > 0 && issues[0].Severity != models.HIGH {
				t.Errorf("expected Severity HIGH, got %v", issues[0].Severity)
			}
		})
	}
}
