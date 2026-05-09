package analyzer

import (
	"config-analyzer/internal/models"
	"config-analyzer/internal/parser"
	"config-analyzer/internal/rules"
	"context"
	"fmt"
)

// Analyzer is the core component for configuration analysis.
type Analyzer struct {
	registry *rules.Registry
}

// New creates a new analyzer with a set of default rules.
func New() *Analyzer {
	return &Analyzer{
		registry: rules.NewRegistry(),
	}
}

// NewWithRegistry creates an analyzer with a custom rules registry.
func NewWithRegistry(registry *rules.Registry) *Analyzer {
	return &Analyzer{
		registry: registry,
	}
}

// Analyze analyzes configuration from raw data.
func (a *Analyzer) Analyze(ctx context.Context, req models.AnalysisRequest) (*models.AnalysisResult, error) {
	type parseResult struct {
		config map[string]interface{}
		err    error
	}
	resultCh := make(chan parseResult, 1)

	go func() {
		cfg, err := parser.Parse(req.Reader, req.FilePath)
		resultCh <- parseResult{config: cfg, err: err}
	}()

	var config map[string]interface{}

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("analysis aborted (timeout/cancel): %w", ctx.Err())
	case res := <-resultCh:
		if res.err != nil {
			return nil, res.err
		}
		config = res.config
	}

	issues := a.registry.CheckAll(config, req.FilePath)

	return &models.AnalysisResult{
		FilePath: req.FilePath,
		Issues:   issues,
	}, nil
}

// Registry returns the rules registry for extension.
func (a *Analyzer) Registry() *rules.Registry {
	return a.registry
}
