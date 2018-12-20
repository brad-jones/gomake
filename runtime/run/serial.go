package run

// Serial takes a list of functions and runs each function one by one
// until either an error is returned or they all complete.
func Serial(fns ...func() error) func() error {
	return func() error {
		for _, fn := range fns {
			if err := fn(); err != nil {
				return err
			}
		}
		return nil
	}
}
