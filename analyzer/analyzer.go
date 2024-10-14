package analyzer

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

// New returns new noifassign analyzer.
func New() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "noifassign",
		Doc:      "check for usage of if with assign statements",
		Run:      run,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			ifStmt, ok := n.(*ast.IfStmt)
			if !ok || ifStmt.Init == nil {
				return true
			}

			assign, ok := ifStmt.Init.(*ast.AssignStmt)
			if !ok {
				return true
			}

			if assign.Tok != token.ASSIGN {
				return true
			}

			pass.Reportf(ifStmt.Pos(), "found 'if' statement with initialization")
			return false
		})
	}
	return nil, nil
}
