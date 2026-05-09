package rules

import (
	"config-analyzer/internal/models"
	"testing"
)

func TestWeakAlgorithmRule_Check(t *testing.T) {
	rule := &WeakAlgorithmRule{}

	tests := []struct {
		name          string
		config        map[string]interface{}
		expectedCount int
		expectMedium  bool
	}{
		{
			name:          "Weak MD5",
			config:        map[string]interface{}{"hash_algorithm": "MD5"},
			expectedCount: 1,
			expectMedium:  false,
		},
		{
			name:          "Weak SHA1",
			config:        map[string]interface{}{"cipher": "sha1"},
			expectedCount: 1,
			expectMedium:  true,
		},
		{
			name:          "Safe AES-256",
			config:        map[string]interface{}{"encryption": "aes-256-gcm"},
			expectedCount: 0,
		},
		{
			name:          "No plaintext allowed",
			config:        map[string]interface{}{"algo": "none"},
			expectedCount: 1,
			expectMedium:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := rule.Check(tt.config, "test.yaml")
			if len(issues) != tt.expectedCount {
				t.Errorf("expected %d issues, got %d", tt.expectedCount, len(issues))
			}

			if tt.expectedCount > 0 {
				expectedSeverity := models.HIGH
				if tt.expectMedium {
					expectedSeverity = models.MEDIUM
				}
				if issues[0].Severity != expectedSeverity {
					t.Errorf("expected Severity %v, got %v", expectedSeverity, issues[0].Severity)
				}
			}
		})
	}
}
