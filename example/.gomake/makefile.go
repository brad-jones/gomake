package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/fatih/color"
	"gopkg.in/brad-jones/gomake.v2/runtime/exec"
	"gopkg.in/brad-jones/gomake.v2/runtime/print"
	"gopkg.in/brad-jones/gomake.v2/runtime/run"
)

// Use can optionally be set to customise the generated usage text.
// It is passed directly to "cobra.Command.Use".
// This is handy if you have built a static binary tool with the generator.
//var Use = "foocmd"

// Short can optionally be set to customise the generated usage text.
// It is passed directly to "cobra.Command.Short".
// This is handy if you have built a static binary tool with the generator.
//var Short = "This is a custom CLI tool"

// Version can optionally be set to customise the generated usage text.
// It is passed directly to "cobra.Command.Version".
//
// If not provided it will default to the current git commit hash of the repo
// that this file is part of, in the event a git repo is not found 0.0.0 will
// be used.
var Version = "0.0.0"

// Validparams shows all valid parameter types
func Validparams(
	qux bool,
	quxArray []bool,
	foo string,
	fooArray []string,
	baz int,
	bazArray []int,
	abc int8,
	xyz int16,
	zxc int32,
	vbn int64,
	asd float32,
	fgh float64,
	jkl uint,
	jklArray []uint,
	qwe uint8,
	rty uint16,
	uio uint32,
	qaz uint64,
	wsx time.Duration,
	wsxArray []time.Duration,
	edc net.IP,
	edcArray []net.IP,
	rfv net.IPMask,
	tgb []byte,
	args ...string,
) error {
	return nil
}

// Sub is an example parent command
func Sub() error {
	fmt.Println("Sub says hi")
	return nil
}

// SubCmd is an example child command
func SubCmd() error {
	fmt.Println("sub cmd says hi")
	return nil
}

// Hyphenated_cmd is an example of a hyphenated-cmd
// Use snake case if your command has many words but should not be considered a
// child command.
func Hyphenated_cmd() error {
	fmt.Println("Hyphenated_cmd says hi")
	return nil
}

// Hyphenated_cmdFoo_bar is an example of a child command of a hyphenated-cmd
// Feel free to mix and match both hyphenated and non hyphenated command names.
// Yeah I know it's a pretty ugly function name but I don't expect this to be
// used as much as sub commands.
func Hyphenated_cmdFoo_bar() error {
	fmt.Println("Hyphenated_cmdFoo_bar says hi")
	return nil
}

// Documentedflag is an example of how to add descriptions against a parameter
// --foo: Is an example parameter
// Flag descriptions must only span a single line, this is not part of
// the --foo description.
func Documentedflag(foo string) error {
	fmt.Println("Documentedflag = " + foo)
	return nil
}

// Shortflag is an example of how to set a shortned alias against a parameter
// --foo, -f: Is an example parameter, with a short alias
func Shortflag(foo string) error {
	fmt.Println("Shortflag = " + foo)
	return nil
}

// Cmd_with_context is an example of a command that accepts a context.
// A default context is passed into any target with a context argument.
// This context will have a timeout if gomake was run with -t, and thus will
// cancel any running functions that the context has been passed to.
func Cmd_with_context(ctx context.Context) error {
	return nil
}

// Noerror is also a valid command
// But discouraged as errors should be handled gracefully instead of panicing.
// Also makes it difficult to use with thing like "run.Serial/Parallel".
func Noerror() {

}

// Cmdwithquotes is a test case to prove gomake can handle quotes in comments.
// For example like "this" and also like 'this'.
func Cmdwithquotes() {

}

// RuntimeExample is an example of using the gomake/runtime library
// This is of course completely optional and what you put inside your
// gomake functions is totally up to you, it's just go afterall.
func RuntimeExample() error {
	return run.Serial(
		print.H1("Start", color.FgRed),
		SubCmd,
		Hyphenated_cmdFoo_bar,
		func() error { return Shortflag("bar") },
		print.H2("These next commands will run asynchronously"),
		run.Parallel(
			Hyphenated_cmd,
			Sub,
			func() error { return Documentedflag("baz") },
			exec.RunPrefixed("google1", "ping", "-c", "4", "8.8.8.8"),
			exec.RunPrefixed("google2", "ping", "-c", "4", "8.8.4.4"),
		),
		print.H1("End", color.FgGreen),
	)()
}
