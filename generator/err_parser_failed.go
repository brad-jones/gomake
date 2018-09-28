package generator

// ErrParserFailed is returned when the golang parser fails.
// This would usually be due to a syntax issue in your makefile
// ".go" source files. See the InnerError for cause.
type ErrParserFailed struct {
	innerError error
}

func (e *ErrParserFailed) Error() string {
	return "gomake: golang parser failed"
}

func (e *ErrParserFailed) InnerError() error {
	return e.innerError
}
