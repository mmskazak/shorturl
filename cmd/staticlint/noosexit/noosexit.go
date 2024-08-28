package noosexit

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

// Analyzer запрещает прямые вызовы os. Exit в основной функции main.
var Analyzer = &analysis.Analyzer{
	Name: "noosexit",
	Doc:  "запрещает прямые вызовы os.Exit в основной функции main",
	Run:  run,
}

// Run функция запускает анализ.
func run(pass *analysis.Pass) (interface{}, error) {
	// Анализируем только файлы из пакета main.
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Проверяем все файлы в пакете
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			// Проверяем функции
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Смотрим чтобы это была функция main
			if fn.Name.Name != "main" {
				continue
			}

			// Проверка вызова os. Exit
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}

				// Проверяем что это вызов os. Exit
				sel, ok := call.Fun.(*ast.SelectorExpr)
				if ok {
					ident, ok := sel.X.(*ast.Ident)
					if ok && ident.Name == "os" && sel.Sel.Name == "Exit" {
						pass.Reportf(call.Pos(), "direct call to os.Exit is not allowed in main package")
					}
				}
				return true
			})
		}
	}
	return nil, nil
}
