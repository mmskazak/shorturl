package noosexit

import (
	"fmt"
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
	var issues []string

	if pass.Pkg.Name() != "main" {
		return issues, nil
	}

	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Name.Name != "main" {
				continue
			}

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
						issue := fmt.Sprintf("Direct call to os.Exit found at position %d", call.Pos())
						issues = append(issues, issue)
						pass.Reportf(call.Pos(), "direct call to os.Exit is not allowed in main package")
					}
				}
				return true
			})
		}
	}

	return issues, nil // Возвращаем список найденных проблем
}
