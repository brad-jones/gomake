package main

import (
	"fmt"
	"os"

	"github.com/brad-jones/gomake/v3/executor"
)

// Injected by ldflags
// Makefile release target does this
// see: https://stackoverflow.com/questions/11354518
var (
	version = "0.0.0"
	commit  = "dev"
	date    = "unknown"
)

func main() {
	// TODO: Figure out what error handling looks like now...

	// As this executable is just a proxy for the underlying task runner we
	// don't want to provide too much custom functionality (cli options/arguments)
	// but it would be nice to be able to inspect the current installed version
	// of this executable.
	if len(os.Args) > 1 && os.Args[1] == "-gmv" {
		fmt.Println("Version: " + version)
		fmt.Println("Commit Ref: " + commit)
		fmt.Println("Build Timestamp: " + date)
		os.Exit(0)
	}

	// Grab the current working dir, the executor uses this to find a makefile.go file
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Call the executor, which locates a makefile.go file, then calls the
	// generator and finally executes the task runner. Caching and other things
	// happen, for more details see the source of the respective packages.
	if err := executor.Execute(cwd, os.Args[1:]...); err != nil {
		panic(err)
	}

	os.Exit(0)
}
