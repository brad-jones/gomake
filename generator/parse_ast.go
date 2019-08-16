package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func parseAST(cwd string) (*token.FileSet, string, *ast.Package, error) {

	fset := token.NewFileSet()

	ignoreGeneratedMakefile := func(file os.FileInfo) bool {
		return file.Name() != "makefile_generated.go"
	}

	pkgs, err := parser.ParseDir(fset, cwd, ignoreGeneratedMakefile, parser.ParseComments)
	if err != nil {
		return nil, "", nil, &ErrParserFailed{ /*innerError: err*/ }
	}

	for k, v := range pkgs {
		if !strings.HasSuffix(k, "_test") {
			if os.Getenv("GOMAKE_PARSE_DEBUG") == "1" {
				ast.Print(fset, v)
			}

			return fset, k, v, nil
		}
	}

	return nil, "", nil, &ErrMainPackageNotFound{ /*innerError: err*/ }
}
