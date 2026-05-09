package rules

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFilePermissionsRule_Check(t *testing.T) {
	// Skip this test on Windows because permissions work differently (ACL vs POSIX)
	if runtime.GOOS == "windows" {
		t.Skip("Skipping file permissions tests on Windows")
	}

	rule := &FilePermissionsRule{}

	// Create a temporary directory for tests
	tempDir := t.TempDir()

	// Create a safe file (0600 - read/write only owner)
	safeFile := filepath.Join(tempDir, "safe.yaml")
	os.WriteFile(safeFile, []byte("test"), 0600)

	// Create an unsafe file (0666 - read/write by everyone)
	unsafeFile := filepath.Join(tempDir, "unsafe.yaml")
	os.WriteFile(unsafeFile, []byte("test"), 0666)

	tests := []struct {
		name          string
		filePath      string
		expectedCount int // Number of found issues
	}{
		{
			name:          "Safe file permissions (0600)",
			filePath:      safeFile,
			expectedCount: 0,
		},
		{
			name:          "Unsafe file permissions (0666)",
			filePath:      unsafeFile,
			expectedCount: 2, // Should find "readable by others" (MEDIUM) and "writable by others" (HIGH)
		},
		{
			name:          "Empty file path",
			filePath:      "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This rule does not use the config map itself, only the file path
			issues := rule.Check(nil, tt.filePath)

			if len(issues) != tt.expectedCount {
				t.Errorf("expected %d issues, got %d", tt.expectedCount, len(issues))
			}
		})
	}
}
