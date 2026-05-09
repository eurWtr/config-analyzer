package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// PlainPasswordRule checks for plaintext passwords or secrets.
type PlainPasswordRule struct{}

// Name returns the rule name
func (r *PlainPasswordRule) Name() string {
	return "plain-password"
}

// sensitiveKeys — keywords that indicate secret data.
var sensitiveKeys = []string{
	"password", "passwd", "pass", "secret", "api_key", "apikey",
	"api-key", "token", "private_key", "private-key", "privatekey",
	"access_key", "access-key", "accesskey", "secret_key", "secret-key",
	"secretkey", "credentials", "auth_token", "auth-token",
}

// Check inspects the configuration and returns found issues
func (r *PlainPasswordRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)

		for _, sensitive := range sensitiveKeys {
			if strings.Contains(lowerKey, sensitive) {
				strVal := fmt.Sprintf("%v", value)
				if strVal != "" && strVal != "<nil>" &&
					!strings.HasPrefix(strVal, "${") &&
					!strings.HasPrefix(strVal, "$") &&
					!strings.HasPrefix(strVal, "vault://") &&
					!strings.HasPrefix(strVal, "env:") &&
					!strings.EqualFold(strVal, "***") {
					issues = append(issues, models.Issue{
						Severity:       models.HIGH,
						Description:    fmt.Sprintf("обнаружен пароль/секрет в открытом виде (ключ: %s)", key),
						Recommendation: "Используйте переменные окружения, vault или другие безопасные способы хранения секретов.",
						Path:           key,
					})
				}
				break
			}
		}
	}

	return issues
}
