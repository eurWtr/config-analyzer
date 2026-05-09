package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// DebugModeRule checks whether debug mode is enabled.
type DebugModeRule struct{}

// Name returns the rule name
func (r *DebugModeRule) Name() string {
	return "debug-mode"
}

// Check inspects the configuration and returns found issues
func (r *DebugModeRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)

		if strings.Contains(lowerKey, "debug") {
			if boolVal, ok := toBool(value); ok && boolVal {
				issues = append(issues, models.Issue{
					Severity:       models.LOW,
					Description:    fmt.Sprintf("debug mode enabled (key: %s)", key),
					Recommendation: "Disable debug mode in production environments.",
					Path:           key,
				})
			}
		}

		if strings.Contains(lowerKey, "level") || strings.Contains(lowerKey, "log_level") || strings.Contains(lowerKey, "loglevel") {
			if strVal, ok := value.(string); ok {
				if strings.EqualFold(strVal, "debug") || strings.EqualFold(strVal, "trace") {
					issues = append(issues, models.Issue{
						Severity:       models.LOW,
						Description:    fmt.Sprintf("logging set to debug level (key: %s, value: %s)", key, strVal),
						Recommendation: "Set a less verbose logging level (info or higher).",
						Path:           key,
					})
				}
			}
		}
	}

	return issues
}

// toBool attempts to coerce a value to boolean
func toBool(v interface{}) (bool, bool) {
	switch val := v.(type) {
	case bool:
		return val, true
	case string:
		lower := strings.ToLower(val)
		if lower == "true" || lower == "yes" || lower == "1" || lower == "on" {
			return true, true
		}
		if lower == "false" || lower == "no" || lower == "0" || lower == "off" {
			return false, true
		}
	case int:
		return val != 0, true
	case float64:
		return val != 0, true
	}
	return false, false
}
