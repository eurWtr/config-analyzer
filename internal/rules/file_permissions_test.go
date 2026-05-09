package rules

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestFilePermissionsRule_Check(t *testing.T) {
	// Пропускаем этот тест на Windows, так как там другая система прав (ACL вместо POSIX)
	if runtime.GOOS == "windows" {
		t.Skip("Пропуск тестирования прав файлов на Windows")
	}

	rule := &FilePermissionsRule{}

	// Создаем временную папку для тестов
	tempDir := t.TempDir()

	// Создаем безопасный файл (600 - чтение/запись только владельцу)
	safeFile := filepath.Join(tempDir, "safe.yaml")
	os.WriteFile(safeFile, []byte("test"), 0600)

	// Создаем опасный файл (666 - чтение/запись всем)
	unsafeFile := filepath.Join(tempDir, "unsafe.yaml")
	os.WriteFile(unsafeFile, []byte("test"), 0666)

	tests := []struct {
		name          string
		filePath      string
		expectedCount int // Количество найденных проблем
	}{
		{
			name:          "Safe file permissions (0600)",
			filePath:      safeFile,
			expectedCount: 0,
		},
		{
			name:          "Unsafe file permissions (0666)",
			filePath:      unsafeFile,
			expectedCount: 2, // Должен найти "чтение всем" (MEDIUM) и "запись всем" (HIGH)
		},
		{
			name:          "Empty file path",
			filePath:      "",
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Это правило не использует саму мапу конфига, только путь к файлу
			issues := rule.Check(nil, tt.filePath)

			if len(issues) != tt.expectedCount {
				t.Errorf("expected %d issues, got %d", tt.expectedCount, len(issues))
			}
		})
	}
}
