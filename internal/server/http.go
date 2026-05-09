package server

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"config-analyzer/internal/analyzer"
	"config-analyzer/internal/models"
)

// HTTPServer предоставляет REST API для анализа конфигураций.
type HTTPServer struct {
	analyzer *analyzer.Analyzer
	addr     string
}

// NewHTTPServer создаёт новый HTTP-сервер.
func NewHTTPServer(addr string, a *analyzer.Analyzer) *HTTPServer {
	return &HTTPServer{
		analyzer: a,
		addr:     addr,
	}
}

// analyzeRequest — структура запроса для REST API.
type analyzeRequest struct {
	Config string `json:"config"`
	Format string `json:"format"`
}

// analyzeResponse — структура ответа REST API.
type analyzeResponse struct {
	Issues []issueResponse `json:"issues"`
	Count  int             `json:"count"`
}

// issueResponse представляет структуру одной проблемы
type issueResponse struct {
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	Recommendation string `json:"recommendation"`
	Path           string `json:"path,omitempty"`
}

// Start запускает HTTP-сервер.
func (s *HTTPServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/analyze", s.handleAnalyze)
	mux.HandleFunc("/health", s.handleHealth)

	slog.Info("HTTP-сервер запущен", "port", s.addr)
	return http.ListenAndServe(s.addr, mux)
}

// handleHealth возвращает состояние сервера
func (s *HTTPServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// handleAnalyze принимает запрос на анализ конфига
func (s *HTTPServer) handleAnalyze(w http.ResponseWriter, r *http.Request) {

	start := time.Now()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var req analyzeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Error("Ошибка парсинга запроса", "error", err)
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	result, err := s.analyzer.Analyze(ctx, models.AnalysisRequest{
		Reader: strings.NewReader(req.Config),
	})

	if err != nil {
		slog.Error("Ошибка анализа", "error", err)
		http.Error(w, "Ошибка: "+err.Error(), http.StatusUnprocessableEntity)
		return
	}

	resp := analyzeResponse{
		Issues: make([]issueResponse, 0, len(result.Issues)),
		Count:  len(result.Issues),
	}

	for _, issue := range result.Issues {
		resp.Issues = append(resp.Issues, issueResponse{
			Severity:       issue.Severity.String(),
			Description:    issue.Description,
			Recommendation: issue.Recommendation,
			Path:           issue.Path,
		})
	}

	slog.Info("Запрос обработан", "duration_ms", time.Since(start).Milliseconds(), "issues", resp.Count)

	w.Header().Set("Content-Type", "application/json")
	if result.HasIssues() {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(resp)
}
