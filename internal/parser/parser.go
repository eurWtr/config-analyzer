package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format определяет формат конфигурационного файла.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatYAML
)

// Parse парсит данные конфигурации в универсальную map-структуру.
// Автоматически определяет формат (JSON или YAML).
func Parse(r io.Reader, filename string) (map[string]interface{}, error) {
	var format Format

	if filename != "" {
		format = DetectFormatByExtension(filename)
	}

	bufReader := bufio.NewReader(r)

	if format == FormatUnknown {
		format = detectFormatFromStream(bufReader)
	}

	switch format {
	case FormatJSON:
		return parseJSON(bufReader)
	case FormatYAML:
		return parseYAML(bufReader)
	default:
		if result, err := parseJSON(bufReader); err == nil {
			return result, nil
		}
		if result, err := parseYAML(bufReader); err == nil {
			return result, nil
		}
		return nil, fmt.Errorf("не удалось распарсить конфигурацию: неизвестный формат")
	}
}

// DetectFormatByExtension определяет формат по расширению файла.
func DetectFormatByExtension(filename string) Format {
	lower := strings.ToLower(filename)
	switch {
	case strings.HasSuffix(lower, ".json"):
		return FormatJSON
	case strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml"):
		return FormatYAML
	default:
		return FormatUnknown
	}
}

// detectFormat определяет формат входных данных
func detectFormatFromStream(r *bufio.Reader) Format {
	peekBytes, err := r.Peek(100)
	if err != nil && err != io.EOF {
		return FormatUnknown
	}

	trimmed := strings.TrimSpace(string(peekBytes))
	if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
		return FormatJSON
	}
	return FormatYAML
}

// parseJSON парсит JSON в map
func parseJSON(r io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка потокового парсинга JSON: %w", err)
	}
	return result, nil
}

// parseYAML парсит YAML в map
func parseYAML(r io.Reader) (map[string]interface{}, error) {
	var result map[string]interface{}
	if err := yaml.NewDecoder(r).Decode(&result); err != nil {
		return nil, fmt.Errorf("ошибка потокового парсинга YAML: %w", err)
	}
	return result, nil
}
