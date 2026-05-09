package rules

import (
	"fmt"
	"strings"

	"config-analyzer/internal/models"
)

// WeakAlgorithmRule checks for the use of deprecated or insecure algorithms.
type WeakAlgorithmRule struct{}

// Name returns the rule name
func (r *WeakAlgorithmRule) Name() string {
	return "weak-algorithm"
}

// weakAlgorithms — list of insecure algorithms with severity levels.
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

// algorithmKeys - list of keywords used to identify hashing and encryption algorithm settings
var algorithmKeys = []string{
	"algorithm", "algo", "cipher", "hash", "digest", "encryption",
	"digest-algorithm", "digest_algorithm", "hash-algorithm", "hash_algorithm",
	"cipher-suite", "cipher_suite", "encryption-method", "encryption_method",
}

// Check inspects the configuration for issues related to hashing and encryption algorithms
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
					Description:    fmt.Sprintf("algorithm is too weak - %s (key: %s)", strings.ToUpper(fmt.Sprintf("%v", value)), key),
					Recommendation: "Replace it with a stronger algorithm (e.g. SHA-256, AES-256).",
					Path:           key,
				})
				break
			}
		}
	}

	return issues
}
