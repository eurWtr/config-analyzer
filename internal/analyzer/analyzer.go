package analyzer

import (
	"config-analyzer/internal/models"
	"config-analyzer/internal/parser"
	"config-analyzer/internal/rules"
	"context"
	"fmt"
)

// Analyzer — основной компонент для анализа конфигураций.
type Analyzer struct {
	registry *rules.Registry
}

// New создаёт новый анализатор с набором правил по умолчанию.
func New() *Analyzer {
	return &Analyzer{
		registry: rules.NewRegistry(),
	}
}

// NewWithRegistry создаёт анализатор с пользовательским реестром правил.
func NewWithRegistry(registry *rules.Registry) *Analyzer {
	return &Analyzer{
		registry: registry,
	}
}

// Analyze анализирует конфигурацию из сырых данных.
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
		return nil, fmt.Errorf("анализ прерван (таймаут/отмена): %w", ctx.Err())
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

// Registry возвращает реестр правил для расширения.
func (a *Analyzer) Registry() *rules.Registry {
	return a.registry
}
