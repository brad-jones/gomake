package main

import (
	"context"
	"fmt"
	"net"
	"time"
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

// Validparms shows all valid parmeter types
func Validparms(
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
	return nil
}

// Hyphenated_cmdFoo_bar is an example of a child command of a hyphenated-cmd
// Feel free to mix and match both hyphenated and non hyphenated command names.
// Yeah I know it's a pretty ugly function name but I don't expect this to be
// used as much as sub commands.
func Hyphenated_cmdFoo_bar() error {
	return nil
}

// Documentedflag is an example of how to add descriptions against a parameter
// --foo: Is an example parameter
// Flag descriptions must only span a single line, this is not part of
// the --foo description.
func Documentedflag(foo string) error {
	return nil
}

// Shortflag is an example of how to set a shortned alias against a parameter
// --foo, -f: Is an example parameter, with a short alias
func Shortflag(foo string) error {
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
func Noerror() {

}

// Cmdwithquotes is a test case to prove gomake can handle quotes in comments.
// For example like "this" and also like 'this'.
func Cmdwithquotes() {

}
