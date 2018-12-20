package run

import (
	"sync"
)

// Parallel takes a list of functions, runs each one in it's own goroutine
// and then either returns the first error it encounters or waits for
// all to finish.
func Parallel(fns ...func() error) func() error {
	return func() error {
		var wg sync.WaitGroup
		var doneCh = make(chan bool)
		var errCh = make(chan error)

		for _, fn := range fns {
			f := fn
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := f(); err != nil {
					errCh <- err
				}
			}()
		}

		go func() { wg.Wait(); doneCh <- true }()

		select {
		case err := <-errCh:
			return err
		case <-doneCh:
			return nil
		}
	}
}
