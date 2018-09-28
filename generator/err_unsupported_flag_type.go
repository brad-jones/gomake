package generator

import "fmt"

// ErrUnsupportedFlagType is returned when a makefile function has used a
// parameter type that this generator and / or cobra does not support.
type ErrUnsupportedFlagType struct {
	innerError       error
	OriginalTypeName string
	IsArray          bool
	MappedTypeName   string
}

func (e *ErrUnsupportedFlagType) Error() string {
	return fmt.Sprintf("gomake: the flag type '%s' is unsupported by gomake and/or cobra", e.MappedTypeName)
}

func (e *ErrUnsupportedFlagType) InnerError() error {
	return e.innerError
}
