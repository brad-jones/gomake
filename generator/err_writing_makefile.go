package generator

// ErrWritingMakefile is returned when `writeGeneratedMakefile` fails.
// See InnerError for more details.
type ErrWritingMakefile struct {
	src []byte
}

func (e *ErrWritingMakefile) Error() string {
	return "gomake: failed to write generated makefile"
}

func (e *ErrWritingMakefile) Source() []byte {
	return e.src
}
