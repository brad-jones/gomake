package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

func parseAST(cwd string) (*token.FileSet, *ast.Package, error) {

	fset := token.NewFileSet()

	ignoreGeneratedMakefile := func(file os.FileInfo) bool {
		return file.Name() != "makefile_generated.go"
	}

	pkgs, err := parser.ParseDir(fset, cwd, ignoreGeneratedMakefile, parser.ParseComments)
	if err != nil {
		return nil, nil, &ErrParserFailed{innerError: err}
	}

	if pkgs["main"] == nil {
		return nil, nil, &ErrMainPackageNotFound{}
	}

	if os.Getenv("GOMAKE_PARSE_DEBUG") == "1" {
		ast.Print(fset, pkgs["main"])
	}

	return fset, pkgs["main"], nil
}
