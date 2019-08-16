package generator

import (
	"github.com/brad-jones/goerr"
)

// Generate parses all ".go" source files in the given working directory, based
// on the resulting abstract syntax tree it will output "makefile_generated.go"
// in the same directory, this contains a main function that executes a
// cobra cli application.
func Generate(cwd string) error {

	_, packageName, nodes, err := parseAST(cwd)
	if err != nil {
		return err
	}

	viewModel, err := buildViewModel(packageName, nodes)
	if err != nil {
		return err
	}

	if err := writeGeneratedMakefile(cwd, viewModel); err != nil {
		return err
	}

	return nil
}

func MustGenerate(cwd string) {
	goerr.Check(Generate(cwd))
}
