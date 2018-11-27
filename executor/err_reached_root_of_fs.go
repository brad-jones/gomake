package executor

// ErrReachedRootOfFs is returned when findGoMakeFolder can not find a
// valid ".gomake" folder after having recursed all the way up to the
// root of filesystem.
type ErrReachedRootOfFs struct {
	innerError error
}

func (e *ErrReachedRootOfFs) Error() string {
	return "gomake: failed to find valid '.gomake' folder, reached root of filesystem"
}

func (e *ErrReachedRootOfFs) InnerError() error {
	return e.innerError
}
