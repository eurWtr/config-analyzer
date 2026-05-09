package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// TLSDisabledRule checks for disabled TLS/SSL.
type TLSDisabledRule struct{}

// Name returns the rule name
func (r *TLSDisabledRule) Name() string {
	return "tls-disabled"
}

// tlsKeys - list of keywords used to identify TLS/SSL related settings
var tlsKeys = []string{
	"tls", "ssl", "https", "tls_verify", "tls-verify", "tlsverify",
	"ssl_verify", "ssl-verify", "sslverify", "verify_ssl", "verify-ssl",
	"verify_tls", "verify-tls", "insecure", "insecure_skip_verify",
	"insecure-skip-verify", "insecureskipverify", "skip_tls_verify",
	"skip-tls-verify",
}

// Check inspects the configuration and returns found issues
func (r *TLSDisabledRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)

		isTLSKey := false
		for _, tk := range tlsKeys {
			if strings.Contains(lowerKey, tk) {
				isTLSKey = true
				break
			}
		}

		if !isTLSKey {
			continue
		}

		// Check keys like "insecure", "skip_verify" — dangerous if true
		if strings.Contains(lowerKey, "insecure") || strings.Contains(lowerKey, "skip") {
			if boolVal, ok := toBool(value); ok && boolVal {
				issues = append(issues, models.Issue{
					Severity:       models.HIGH,
					Description:    fmt.Sprintf("отключена TLS-верификация (ключ: %s)", key),
					Recommendation: "Включите проверку TLS-сертификатов для защиты от MITM-атак.",
					Path:           key,
				})
			}
			continue
		}

		// Check keys like "tls.enabled", "ssl" — dangerous if false
		if strings.Contains(lowerKey, "enabled") || strings.HasSuffix(lowerKey, "tls") ||
			strings.HasSuffix(lowerKey, "ssl") || strings.HasSuffix(lowerKey, "https") {
			if boolVal, ok := toBool(value); ok && !boolVal {
				issues = append(issues, models.Issue{
					Severity:       models.HIGH,
					Description:    fmt.Sprintf("TLS/SSL отключен (ключ: %s)", key),
					Recommendation: "Включите TLS для шифрования трафика.",
					Path:           key,
				})
			}
		}
	}

	return issues
}
