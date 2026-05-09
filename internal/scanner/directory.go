package scanner

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"config-analyzer/internal/analyzer"
	"config-analyzer/internal/models"
)

// supportedExtensions — file extensions that will be analyzed.
var supportedExtensions = map[string]bool{
	".json": true,
	".yaml": true,
	".yml":  true,
}

// ScanDirectory recursively scans a directory and analyzes configuration files.
func ScanDirectory(dir string, a *analyzer.Analyzer) ([]models.AnalysisResult, error) {
	var results []models.AnalysisResult

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			slog.Warn("Ошибка доступа к пути", "path", path, "error", err)
			return nil
		}

		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !supportedExtensions[ext] {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			slog.Warn("Не удалось открыть файл", "file", path, "error", err)
			return nil
		}
		defer file.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		result, err := a.Analyze(ctx, models.AnalysisRequest{
			Reader:   file,
			FilePath: path,
		})

		if err != nil {
			slog.Warn("Ошибка анализа файла", "file", path, "error", err)
			return nil
		}

		if result.HasIssues() {
			results = append(results, *result)
		}
		return nil
	})

	return results, err
}
