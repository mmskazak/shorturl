package staticlint

import (
	"encoding/json"
	"fmt"
	"github.com/alexkohler/nakedret"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
	"log"
	"mmskazak/shorturl/cmd/staticlint/noosexit"
	"os"
)

// Допустимое количество строк в функции для возврата голого ответа
const countLinesNakedFunc = 25

func main() {
	// Соберем все стандартные анализаторы
	analyzers := []*analysis.Analyzer{
		buildtag.Analyzer,
		copylock.Analyzer,
		fieldalignment.Analyzer,
		httpresponse.Analyzer,
		loopclosure.Analyzer,
		nilfunc.Analyzer,
		printf.Analyzer,
		shift.Analyzer,
		sortslice.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		// ...другие анализаторы по необходимости
	}

	// определяем map подключаемых правил
	//checks := map[string]bool{
	//	"S1001":  true,
	//	"ST1001": true,
	//	"QF1002": true,
	//}

	// Загружаем конфигурацию
	checks, err := loadConfig("config.json")
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}

	var myChecks []*analysis.Analyzer

	// Анализаторы из пакета staticcheck
	for _, v := range staticcheck.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	for _, v := range simple.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	for _, v := range stylecheck.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	for _, v := range quickfix.Analyzers {
		if checks[v.Analyzer.Name] {
			myChecks = append(myChecks, v.Analyzer)
		}
	}

	// Добавим внешний анализатор 1
	analyzers = append(analyzers, errcheck.Analyzer)

	// Добавим внешний анализатор 2 (проверка на голые возвраты)
	analyzers = append(analyzers, nakedret.NakedReturnAnalyzer(countLinesNakedFunc))

	// Добавим собственный анализатор
	analyzers = append(analyzers, noosexit.Analyzer)

	multichecker.Main(analyzers...)
}

type Config struct {
	Checks map[string]bool `json:"checks"`
}

// Считываем конфигурацию из файла
func loadConfig(path string) (map[string]bool, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл конфигурации: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("Error file config not close: %v", err)
		}
	}(file)

	var config Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("не удалось декодировать файл конфигурации: %v", err)
	}

	return config.Checks, nil
}
