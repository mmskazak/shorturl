package noosexit

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData() // Путь до папки с тестовыми данными

	// Запуск теста анализатора
	analysistest.Run(t, testdata, Analyzer, "a")
}
