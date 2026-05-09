package models

import (
	"fmt"
	"io"
)

// Severity определяет уровень серьёзности проблемы.
type Severity int

const (
	LOW Severity = iota
	MEDIUM
	HIGH
)

func (s Severity) String() string {
	switch s {
	case LOW:
		return "LOW"
	case MEDIUM:
		return "MEDIUM"
	case HIGH:
		return "HIGH"
	default:
		return "UNKNOWN"
	}
}

// Issue представляет найденную проблему в конфигурации.
type Issue struct {
	Severity       Severity `json:"severity"`
	Description    string   `json:"description"`
	Recommendation string   `json:"recommendation"`
	Path           string   `json:"path,omitempty"` // путь к ключу в конфиге
}

// String возвращает текстовое представление проблемы
func (i Issue) String() string {
	s := fmt.Sprintf("%s: %s. %s", i.Severity, i.Description, i.Recommendation)
	if i.Path != "" {
		s = fmt.Sprintf("[%s] %s", i.Path, s)
	}
	return s
}

// AnalysisRequest содержит данные для анализа.
type AnalysisRequest struct {
	Reader   io.Reader
	FilePath string
}

// AnalysisResult содержит результат анализа.
type AnalysisResult struct {
	FilePath string  `json:"file_path,omitempty"`
	Issues   []Issue `json:"issues"`
}

// HasIssues проверяет наличие проблем
func (r AnalysisResult) HasIssues() bool {
	return len(r.Issues) > 0
}
