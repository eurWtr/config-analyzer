package rules

import (
	"fmt"
	"os"

	"config-analyzer/internal/models"
)

// FilePermissionsRule checks file permissions of the configuration file.
type FilePermissionsRule struct{}

// Name returns the rule name
func (r *FilePermissionsRule) Name() string {
	return "file-permissions"
}

// Check inspects the file path and returns found issues
func (r *FilePermissionsRule) Check(_ map[string]interface{}, filePath string) []models.Issue {
	var issues []models.Issue

	if filePath == "" {
		return issues
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return issues
	}

	mode := info.Mode().Perm()

	// Проверяем, доступен ли файл для чтения всем (other)
	if mode&0o004 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.MEDIUM,
			Description:    fmt.Sprintf("конфигурационный файл доступен для чтения всем пользователям (права: %o)", mode),
			Recommendation: "Ограничьте права доступа: chmod 600 или chmod 640.",
			Path:           filePath,
		})
	}

	// Проверяем, доступен ли файл для записи всем (other)
	if mode&0o002 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.HIGH,
			Description:    fmt.Sprintf("конфигурационный файл доступен для записи всем пользователям (права: %o)", mode),
			Recommendation: "Немедленно ограничьте права доступа: chmod 600.",
			Path:           filePath,
		})
	}

	// Проверяем, доступен ли файл для записи группе
	if mode&0o020 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.LOW,
			Description:    fmt.Sprintf("конфигурационный файл доступен для записи группе (права: %o)", mode),
			Recommendation: "Рассмотрите ограничение прав: chmod 640 или chmod 600.",
			Path:           filePath,
		})
	}

	return issues
}
