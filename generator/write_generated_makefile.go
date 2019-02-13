package generator

import (
	"bytes"
	"go/format"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"gopkg.in/brad-jones/gomake.v2/resources"
)

func writeGeneratedMakefile(cwd string, viewModel *tplViewModel) error {

	tpl, err := template.New("generated").Parse(resources.Typed.MakefileGeneratedGotmpl.String())
	if err != nil {
		return &ErrWritingMakefile{innerError: err}
	}

	var unFormatted bytes.Buffer
	if err := tpl.Execute(&unFormatted, viewModel); err != nil {
		return &ErrWritingMakefile{innerError: err}
	}

	formatted, err := format.Source(unFormatted.Bytes())
	if err != nil {
		return &ErrWritingMakefile{innerError: err}
	}

	if err := ioutil.WriteFile(filepath.Join(cwd, "makefile_generated.go"), formatted, 0644); err != nil {
		return &ErrWritingMakefile{innerError: err}
	}

	return nil
}
