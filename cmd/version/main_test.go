package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

// Функция для очистки отступов и пробелов
func normalizeWhitespace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func TestMainFunction(t *testing.T) {
	// Сохранить текущие переменные окружения
	oldBuildVersion := os.Getenv("BUILD_VERSION")
	oldBuildCommit := os.Getenv("BUILD_COMMIT")

	// Установить тестовые переменные окружения
	os.Setenv("BUILD_VERSION", "v1.0.0")
	os.Setenv("BUILD_COMMIT", "abcdef1234")

	// Запустить функцию main
	main()

	// Проверить, что файл создан и содержит правильный код
	data, err := os.ReadFile("version_gen.go")
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	expectedBuildDate := time.Now().Format(time.DateTime)
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
	os.Remove("version_gen.go")

	// Восстановить старые переменные окружения
	os.Setenv("BUILD_VERSION", oldBuildVersion)
	os.Setenv("BUILD_COMMIT", oldBuildCommit)
}
