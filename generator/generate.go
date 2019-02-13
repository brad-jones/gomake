package generator

// Generate parses all ".go" source files in the given working directory, based
// on the resulting abstract syntax tree it will output "makefile_generated.go"
// in the same directory, this contains a main function that executes a
// cobra cli application.
func Generate(cwd string) error {

	_, nodes, err := parseAST(cwd)
	if err != nil {
		return err
	}

	viewModel, err := buildViewModel(nodes)
	if err != nil {
		return err
	}

	h, err := CacheHashGen(cwd)
	if err != nil {
		return err
	}
	viewModel.CacheHash = h

	if err := writeGeneratedMakefile(cwd, viewModel); err != nil {
		return err
	}

	return nil
}
