package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Простой хендлер для тестирования
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "test response"}`)) //nolint:gosec,errcheck //можно пренебречь
}

func TestGzipMiddleware(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(testHandler))

	t.Run("Without Gzip", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", http.NoBody)
		rec := httptest.NewRecorder()

		handler.ServeHTTP(rec, req)

		res := rec.Result()
		defer res.Body.Close() //nolint:errcheck //можно пренебречь

		if res.Header.Get("Content-Encoding") != "" {
			t.Errorf("expected no Content-Encoding, got %s", res.Header.Get("Content-Encoding"))
		}

		body, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		expected := `{"message": "test response"}`
		if strings.TrimSpace(string(body)) != expected {
			t.Errorf("expected %s, got %s", expected, body)
		}
	})
}

func TestGzipMiddleware_RequestWithGzip(t *testing.T) {
	handler := GzipMiddleware(http.HandlerFunc(testHandler))

	// Создаем сжатый gzip запрос
	var requestBody strings.Builder
	gzipWriter := gzip.NewWriter(&requestBody)
	gzipWriter.Write([]byte(`{"message": "test request"}`)) //nolint:gosec,errcheck //можно пренебречь
	gzipWriter.Close()                                      //nolint:gosec,errcheck //можно пренебречь

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(requestBody.String()))
	req.Header.Set("Content-Encoding", "gzip")

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close() //nolint:errcheck //можно пренебречь

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.Status)
	}

	// Проверка ответа от обработчика
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := `{"message": "test response"}`
	if strings.TrimSpace(string(body)) != expected {
		t.Errorf("expected %s, got %s", expected, body)
	}
}
