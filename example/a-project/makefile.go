package makeprja

import (
	"fmt"

	"github.com/brad-jones/goasync/task"
)

// This file can be used as a standalone task runner or you could import it in
// some other parent task runner project and use it's exported functions programatically.

func BuildProjectA() *task.Task {
	return task.New(func(t *task.Internal) {
		fmt.Println("Building Project A...")
	})
}
