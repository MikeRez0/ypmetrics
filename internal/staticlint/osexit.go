package staticlint

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

var OsExitCheck = &analysis.Analyzer{
	Name: "osexit",
	Doc:  "check for call os.Exit() in main func",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var err error

	for _, file := range pass.Files {
		if file.Name.Name != "main" {
			continue
		}

		parents := make([]*ast.Node, 0)
		ast.Inspect(file, func(node ast.Node) bool {
			if node == nil && len(parents) > 0 {
				parents = parents[:len(parents)-1]
			} else if len(parents) >= 3 {
				parents = append(parents, &node)
			}

			switch x := node.(type) {
			case *ast.File:
				if x.Name.Name == "main" {
					parents = append(parents, &node)
					return true
				}
			case *ast.FuncDecl:
				if x.Name.Name == "main" {
					parents = append(parents, &node)
					return true
				}
			case *ast.BlockStmt:
				parents = append(parents, &node)
				return true
			case *ast.CallExpr:
				if s, ok := x.Fun.(*ast.SelectorExpr); ok {
					if s.Sel.Name == "Exit" {
						if p, ok := s.X.(*ast.Ident); ok {
							if p.Name == "os" {
								pass.Reportf(x.Pos(), "os.Exit call from main is not recommended")
								return false
							}
						}
					}
				}
			}
			if len(parents) >= 3 {
				return true
			} else {
				return false
			}
		})
	}

	return nil, err
}
