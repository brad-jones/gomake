package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/brad-jones/goerr"
	"github.com/brad-jones/gomake/v3/executor"
	"github.com/brad-jones/gomake/v3/generator"
	"github.com/fatih/color"
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
	// Recover from all panics here to make error handling uniform and maintainable.
	defer goerr.Handle(func(err error) {
		red := color.New(color.FgRed, color.Bold, color.Underline)
		red.Println("gomake has encountered an error!")
		fmt.Println(err)
		fmt.Println()
		if os.Getenv("GOMAKE_DEBUG") == "1" {
			red.Println("debug stack trace")
			if stackTrace, _ := goerr.Trace(err); stackTrace != "" {
				fmt.Println(stackTrace)
			} else {
				debug.PrintStack()
			}
		}
		os.Exit(1)
	})

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

	// Call the executor, which locates a makefile.go file, then calls the
	// generator and finally executes the task runner. Caching and other things
	// happen, for more details see the source of the respective packages.
	executor.New(generator.New()).MustExecute(os.Args[1:]...)
}
