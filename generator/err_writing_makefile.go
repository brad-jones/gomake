package generator

// ErrWritingMakefile is returned when `writeGeneratedMakefile` fails.
// See InnerError for more details.
type ErrWritingMakefile struct {
	innerError error
}

func (e *ErrWritingMakefile) Error() string {
	return "gomake: failed to write generated makefile"
}

func (e *ErrWritingMakefile) InnerError() error {
	return e.innerError
}
