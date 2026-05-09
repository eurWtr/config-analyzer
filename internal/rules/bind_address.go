package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// BindAddressRule проверяет привязку к 0.0.0.0.
type BindAddressRule struct{}

// Name возвращает имя правила
func (r *BindAddressRule) Name() string {
	return "bind-address"
}

// addressKeys список ключевых слов для поиска нужных полей
var addressKeys = []string{
	"host", "bind", "address", "addr", "listen", "bind_address",
	"bind-address", "listen_address", "listen-address", "server",
}

// Check проверяет конфигурацию и возвращает найденные проблемы
func (r *BindAddressRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)
		strVal := fmt.Sprintf("%v", value)

		isAddressKey := false
		for _, ak := range addressKeys {
			if strings.Contains(lowerKey, ak) {
				isAddressKey = true
				break
			}
		}

		if !isAddressKey {
			continue
		}

		if strings.Contains(strVal, "0.0.0.0") {
			issues = append(issues, models.Issue{
				Severity:       models.MEDIUM,
				Description:    fmt.Sprintf("сервис привязан к 0.0.0.0 (ключ: %s), доступен на всех интерфейсах", key),
				Recommendation: "Ограничьте привязку конкретным интерфейсом (например, 127.0.0.1) или используйте firewall.",
				Path:           key,
			})
		}
	}

	return issues
}
