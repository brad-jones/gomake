package makefile

// This is just another file that is included in the main task runner,
// changes to this file should trigger a rebuild.

import (
	makeprja "github.com/brad-jones/gomake/v3/example/a-project"
)

func BuildProjectA() {
	makeprja.BuildProjectA()
}
