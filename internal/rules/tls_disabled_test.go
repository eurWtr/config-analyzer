package rules

import (
	"config-analyzer/internal/models"
	"testing"
)

func TestTLSDisabledRule_Check(t *testing.T) {
	rule := &TLSDisabledRule{}

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectedCount int
	}{
		{
			name:          "TLS explicitly disabled",
			config:        map[string]interface{}{"tls": map[string]interface{}{"enabled": false}},
			expectedCount: 1,
		},
		{
			name:          "Insecure skip verify",
			config:        map[string]interface{}{"insecure_skip_verify": true},
			expectedCount: 1,
		},
		{
			name: "Safe TLS settings",
			config: map[string]interface{}{
				"tls_enabled": true,
				"verify_ssl":  true,
			},
			expectedCount: 0,
		},
		{
			name:          "String boolean representations",
			config:        map[string]interface{}{"ssl": "off"},
			expectedCount: 1,
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
