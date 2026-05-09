package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// DebugModeRule проверяет, включен ли debug-режим.
type DebugModeRule struct{}

// Name возвращает имя правила
func (r *DebugModeRule) Name() string {
	return "debug-mode"
}

// Check проверяет конфигурацию и возвращает найденные проблемы
func (r *DebugModeRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)

		if strings.Contains(lowerKey, "debug") {
			if boolVal, ok := toBool(value); ok && boolVal {
				issues = append(issues, models.Issue{
					Severity:       models.LOW,
					Description:    fmt.Sprintf("включен debug-режим (ключ: %s)", key),
					Recommendation: "Отключите debug-режим в production-окружении.",
					Path:           key,
				})
			}
		}

		if strings.Contains(lowerKey, "level") || strings.Contains(lowerKey, "log_level") || strings.Contains(lowerKey, "loglevel") {
			if strVal, ok := value.(string); ok {
				if strings.EqualFold(strVal, "debug") || strings.EqualFold(strVal, "trace") {
					issues = append(issues, models.Issue{
						Severity:       models.LOW,
						Description:    fmt.Sprintf("логирование в debug-режиме (ключ: %s, значение: %s)", key, strVal),
						Recommendation: "Поменяйте режим на более избирательный (info+).",
						Path:           key,
					})
				}
			}
		}
	}

	return issues
}

// toBool пытается привести значение к логическому типу
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
