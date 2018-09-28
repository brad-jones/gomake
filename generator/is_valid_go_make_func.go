package generator

import (
	"go/ast"
	"strings"
)

func isValidGoMakeFunc(funcDecl *ast.FuncDecl) bool {
	firstFuncChar := string(funcDecl.Name.Name[0])
	if strings.ToUpper(firstFuncChar) == firstFuncChar {
		if funcDecl.Type.Results == nil {
			return true
		}
		if len(funcDecl.Type.Results.List) == 1 {
			if r, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident); ok {
				if r.Name == "error" {
					return true
				}
			}
		}
	}
	return false
}
