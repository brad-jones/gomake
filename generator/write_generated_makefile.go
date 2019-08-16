package generator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/brad-jones/gomake/v3/resources"
)

func writeGeneratedMakefile(cwd string, viewModel *tplViewModel) error {

	tpl, err := template.New("generated").Parse(resources.Typed.MakefileGeneratedGotmpl.String())
	if err != nil {
		return &ErrWritingMakefile{ /* innerError: err */ }
	}

	var unFormatted bytes.Buffer
	if err := tpl.Execute(&unFormatted, viewModel); err != nil {
		return &ErrWritingMakefile{ /* innerError: err */ }
	}

	formatted, err := format.Source(unFormatted.Bytes())
	if err != nil {
		return &ErrWritingMakefile{
			/* innerError: err, */
			src: unFormatted.Bytes(),
		}
	}

	os.MkdirAll(filepath.Join(cwd, ".gomake"), 0744)

	if err := ioutil.WriteFile(filepath.Join(cwd, ".gomake", "makefile_generated.go"), formatted, 0644); err != nil {
		panic(err)
		return &ErrWritingMakefile{ /* innerError: err */ }
	}

	return nil
}
