package generator

// ErrVariadicMustBeString is returned when a makefile function has used a
// variadic parameter that is not of type "string".
type ErrVariadicMustBeString struct {
}

func (e *ErrVariadicMustBeString) Error() string {
	return "gomake: variadic parameters must of type 'string'"
}
