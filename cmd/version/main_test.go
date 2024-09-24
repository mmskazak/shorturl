package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Функция для очистки отступов и пробелов.
func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func TestMainFunction(t *testing.T) {
	// Установить тестовые переменные окружения
	t.Setenv("BUILD_VERSION", "v1.0.0")
	t.Setenv("BUILD_COMMIT", "abcdef1234")
	expectedBuildDate := time.Now().Format(time.DateTime)
	t.Setenv("BUILD_DATE", expectedBuildDate)

	// Запустить функцию main
	main()

	// Проверить, что файл создан и содержит правильный код
	data, err := os.ReadFile("version_gen.go")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expectedCode := fmt.Sprintf(`package main

var (
    BuildVersion = "v1.0.0"
    BuildDate    = "%s"
    BuildCommit  = "abcdef1234"
)
`, expectedBuildDate)

	// Нормализуем пробелы и отступы для сравнения
	if normalizeWhitespace(string(data)) != normalizeWhitespace(expectedCode) {
		t.Errorf("File content is incorrect. Got:\n%s\nExpected:\n%s", data, expectedCode)
	}

	// Удалить файл после тестирования
	err = os.Remove("version_gen.go")
	require.NoError(t, err)
}
