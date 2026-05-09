package models

import (
	"fmt"
	"io"
)

// Severity defines the severity level of an issue.
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

// Issue represents a found configuration issue.
type Issue struct {
	Severity       Severity `json:"severity"`
	Description    string   `json:"description"`
	Recommendation string   `json:"recommendation"`
	Path           string   `json:"path,omitempty"` // path to the key in the config
}

// String returns a textual representation of the issue
func (i Issue) String() string {
	s := fmt.Sprintf("%s: %s. %s", i.Severity, i.Description, i.Recommendation)
	if i.Path != "" {
		s = fmt.Sprintf("[%s] %s", i.Path, s)
	}
	return s
}

// AnalysisRequest holds data for analysis.
type AnalysisRequest struct {
	Reader   io.Reader
	FilePath string
}

// AnalysisResult contains the analysis result.
type AnalysisResult struct {
	FilePath string  `json:"file_path,omitempty"`
	Issues   []Issue `json:"issues"`
}

// HasIssues checks whether any issues were found
func (r AnalysisResult) HasIssues() bool {
	return len(r.Issues) > 0
}
