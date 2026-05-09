package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// WeakAlgorithmRule проверяет использование устаревших или небезопасных алгоритмов.
type WeakAlgorithmRule struct{}

// Name возвращает имя правила
func (r *WeakAlgorithmRule) Name() string {
	return "weak-algorithm"
}

// weakAlgorithms — список небезопасных алгоритмов с уровнем серьёзности.
var weakAlgorithms = map[string]models.Severity{
	"md5":       models.HIGH,
	"md4":       models.HIGH,
	"sha1":      models.MEDIUM,
	"sha-1":     models.MEDIUM,
	"des":       models.HIGH,
	"3des":      models.HIGH,
	"rc4":       models.HIGH,
	"rc2":       models.HIGH,
	"blowfish":  models.MEDIUM,
	"none":      models.HIGH,
	"plaintext": models.HIGH,
}

// algorithmKeys - список ключевых слов, используемых для идентификации алгоритмов хеширования и шифрования
var algorithmKeys = []string{
	"algorithm", "algo", "cipher", "hash", "digest", "encryption",
	"digest-algorithm", "digest_algorithm", "hash-algorithm", "hash_algorithm",
	"cipher-suite", "cipher_suite", "encryption-method", "encryption_method",
}

// Check проверяет конфигурацию на проблемы, связанные с алгоритмами хеширования и шифрования
func (r *WeakAlgorithmRule) Check(config map[string]interface{}, _ string) []models.Issue {
	var issues []models.Issue
	flat := flatten("", config)

	for key, value := range flat {
		lowerKey := strings.ToLower(key)

		isAlgoKey := false
		for _, ak := range algorithmKeys {
			if strings.Contains(lowerKey, ak) {
				isAlgoKey = true
				break
			}
		}

		if !isAlgoKey {
			continue
		}

		strVal := strings.ToLower(fmt.Sprintf("%v", value))

		for weakAlgo, severity := range weakAlgorithms {
			if strings.Contains(strVal, weakAlgo) {
				issues = append(issues, models.Issue{
					Severity:       severity,
					Description:    fmt.Sprintf("слишком слабый алгоритм - %s (ключ: %s)", strings.ToUpper(fmt.Sprintf("%v", value)), key),
					Recommendation: "Замените его на более безопасный (например, SHA-256, AES-256).",
					Path:           key,
				})
				break
			}
		}
	}

	return issues
}
