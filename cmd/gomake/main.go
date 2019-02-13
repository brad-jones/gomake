package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
	"gopkg.in/brad-jones/gomake.v2/executor"
)

// Injected by ldflags
// Makefile release target does this
// see: https://stackoverflow.com/questions/11354518
var (
	version = "0.0.0"
	commit  = "dev"
	date    = "unknown"
)

// Use text that stands out to report errors
var red = color.New(color.FgRed, color.Bold, color.Underline)

// The gomake errors have a notion of an "inner error", similar to other
// languages such C#. This allows errors to provide additional context.
type errWithInnerErr interface {
	InnerError() error
}

// printInnerError recursively prints inner error messages
func printInnerErrors(err interface{}) {
	if err, ok := err.(errWithInnerErr); ok {
		if innerErr := err.InnerError(); innerErr != nil {
			red.Println("inner error")
			fmt.Println(innerErr)
			fmt.Println()
			printInnerErrors(innerErr)
		}
	}
}

// handleAllErrors recovers from all panics to make error handling
// uniform and maintainable, it then exits with a code of 1.
func handleAllErrors() {
	if err := recover(); err != nil {
		// Print the top level error message
		red.Println("gomake has encountered an error!")
		fmt.Println(err)
		fmt.Println()

		// Print any inner errors recursivly
		printInnerErrors(err)

		// Optionally print a stack trace
		if os.Getenv("GOMAKE_DEBUG") == "1" {
			red.Println("debug stack trace")
			debug.PrintStack()
			fmt.Println()
		}

		os.Exit(1)
	}
}

func main() {

	defer handleAllErrors()

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

	// Grab the current working dir, the executor uses this to find a .gomake folder
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Call the executor, which locates a .gomake folder, then calls the
	// generator and finally executes the task runner. Caching and other things
	// happen, for more details see the source of the respective packages.
	if err := executor.Execute(cwd, os.Args[1:]...); err != nil {
		panic(err)
	}

	os.Exit(0)
}
