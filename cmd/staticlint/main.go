package staticlint

import (
	"github.com/alexkohler/nakedret"
	"github.com/kisielk/errcheck/errcheck"
	"golang.org/x/tools/go/analysis"
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
	"mmskazak/shorturl/cmd/staticlint/noosexit"
	"strings"

	"golang.org/x/tools/go/analysis/multichecker"
	"honnef.co/go/tools/staticcheck"
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

	// Анализаторы из пакета staticcheck
	for _, v := range staticcheck.Analyzers {
		// Все анализаторы класса SA пакета staticcheck.io
		if strings.HasPrefix(v.Analyzer.Name, "SA") {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	// определяем map подключаемых правил
	checks := map[string]bool{
		"S1001":  true,
		"ST1001": true,
		"QF1002": true,
	}
	var myChecks []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// добавляем в массив нужные проверки
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
