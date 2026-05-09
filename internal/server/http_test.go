package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"config-analyzer/internal/analyzer"
)

func TestHTTPServer_handleAnalyze(t *testing.T) {
	// 1. Создаем реальный инстанс анализатора
	a := analyzer.New()
	srv := NewHTTPServer(":8080", a)

	// 2. Формируем тело тестового запроса (JSON с "плохим" конфигом)
	reqBody := analyzeRequest{
		Config: `{"db": {"password": "123"}}`,
	}
	bodyBytes, _ := json.Marshal(reqBody)

	// 3. Создаем фейковый HTTP-запрос
	req, err := http.NewRequest(http.MethodPost, "/analyze", bytes.NewReader(bodyBytes))
	if err != nil {
		t.Fatal(err)
	}

	// 4. Создаем Recorder, который запишет ответ сервера
	rr := httptest.NewRecorder()

	// 5. Вызываем хэндлер напрямую!
	srv.handleAnalyze(rr, req)

	// 6. Проверяем HTTP статус
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// 7. Проверяем структуру ответа
	var resp analyzeResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}

	// Ожидаем, что сработает правило PlainPasswordRule
	if resp.Count != 1 {
		t.Errorf("expected 1 issue, got %d", resp.Count)
	}

	if len(resp.Issues) > 0 && resp.Issues[0].Severity != "HIGH" {
		t.Errorf("expected HIGH severity, got %s", resp.Issues[0].Severity)
	}
}

func TestHTTPServer_WrongMethod(t *testing.T) {
	srv := NewHTTPServer(":8080", analyzer.New())

	// Отправляем GET вместо POST
	req, _ := http.NewRequest(http.MethodGet, "/analyze", nil)
	rr := httptest.NewRecorder()

	srv.handleAnalyze(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405 StatusMethodNotAllowed, got %v", rr.Code)
	}
}
