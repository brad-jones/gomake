package makefile

// This is just another file that is included in the main task runner,
// changes to this file should trigger a rebuild.

import (
	"github.com/brad-jones/goasync/task"
	makeprja "github.com/brad-jones/gomake/v3/example/a-project"
)

func BuildProjectA() *task.Task {
	return makeprja.BuildProjectA()
}
