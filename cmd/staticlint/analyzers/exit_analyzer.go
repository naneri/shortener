package analyzers

import (
	"fmt"
	"go/ast"
	"golang.org/x/tools/go/analysis"
)

var ExitAnalyzer = &analysis.Analyzer{
	Name: "exitcheck",
	Doc:  "check for os.Exit() usage",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	// реализация будет ниже
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.CallExpr:
				switch node := x.Fun.(type) {
				case *ast.SelectorExpr:
					//fmt.Println(node.info.Uses[pkgID].(*types.PkgName).Imported().Path())
					res := fmt.Sprintf("%s.%s", node.X, node.Sel.Name)
					if f.Name.Name == "main" && res == "os.Exit" {
						fmt.Printf("%s: avoid using os.Exit() in the main function \n", pass.Fset.Position(node.Pos()))
					}
				}
			}
			return true
		})
	}

	return nil, nil
}
