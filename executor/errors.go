package executor

import (
	"github.com/go-errors/errors"
)

// ErrReachedRootOfFs is returned when the "findGomakeFolder" method recursed
// all the way to the root of the filesystem and failed to find a valid "makefile.go"
var ErrReachedRootOfFs = errors.Errorf("gomake: failed to find a valid 'makefile.go', reached root of filesystem")
