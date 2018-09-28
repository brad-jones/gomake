package generator

// ErrVariadicMustBeString is returned when a makefile function has used a
// variadic parameter that is not of type "string".
type ErrVariadicMustBeString struct {
	innerError error
}

func (e *ErrVariadicMustBeString) Error() string {
	return "gomake: variadic parameters must of type 'string'"
}

func (e *ErrVariadicMustBeString) InnerError() error {
	return e.innerError
}
