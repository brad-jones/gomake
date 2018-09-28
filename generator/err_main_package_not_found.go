package generator

// ErrMainPackageNotFound is returned when a folder of ".go" source files have
// been provided that are not part of a "main" package. Or there simply are no
// source files. All direct children of a ".gomake" folder must be part of the
// "main" package.
type ErrMainPackageNotFound struct {
	innerError error
}

func (e *ErrMainPackageNotFound) Error() string {
	return "gomake: main package not found in parsed ast"
}

func (e *ErrMainPackageNotFound) InnerError() error {
	return e.innerError
}
