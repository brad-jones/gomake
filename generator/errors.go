package generator

import (
	"fmt"

	"github.com/go-errors/errors"
)

// ErrReachedRootOfFs is returned when the "findGoModSumFiles" method recursed
// all the way to the root of the filesystem and failed to find a valid "go.mod"
var ErrReachedRootOfFs = errors.Errorf("gomake: failed to find a valid 'go.mod', reached root of filesystem")

// ErrNoPackageNotFound is returned when the golang parser fails to find any
// package where the "makefile.go" is found.
var ErrNoPackageNotFound = errors.Errorf("gomake: no package found")

// ErrVariadicMustBeString is returned when a makefile function has used a
// variadic parameter that is not of type "string".
var ErrVariadicMustBeString = errors.Errorf("gomake: variadic parameters must be of type 'string'")

// ErrUnsupportedFlagType is returned when a makefile function has used a
// parameter type that this generator and / or cobra does not support.
type ErrUnsupportedFlagType struct {
	OriginalTypeName string
	IsArray          bool
	MappedTypeName   string
}

func (e *ErrUnsupportedFlagType) Error() string {
	return fmt.Sprintf("gomake: the flag type '%s' is unsupported by gomake and/or cobra", e.MappedTypeName)
}
