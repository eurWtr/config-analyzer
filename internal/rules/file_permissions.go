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

	// Check if the file is readable by others
	if mode&0o004 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.MEDIUM,
			Description:    fmt.Sprintf("configuration file is readable by others (mode: %o)", mode),
			Recommendation: "Restrict permissions: chmod 600 or chmod 640.",
			Path:           filePath,
		})
	}

	// Check if the file is writable by others
	if mode&0o002 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.HIGH,
			Description:    fmt.Sprintf("configuration file is writable by others (mode: %o)", mode),
			Recommendation: "Immediately restrict permissions: chmod 600.",
			Path:           filePath,
		})
	}

	// Check if the file is writable by group
	if mode&0o020 != 0 {
		issues = append(issues, models.Issue{
			Severity:       models.LOW,
			Description:    fmt.Sprintf("configuration file is writable by group (mode: %o)", mode),
			Recommendation: "Consider restricting permissions: chmod 640 or chmod 600.",
			Path:           filePath,
		})
	}

	return issues
}
