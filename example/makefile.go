package makefile

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/brad-jones/goasync/task"
	"github.com/brad-jones/gomake/v3/deps"
	errors2 "github.com/go-errors/errors"
	errors1 "github.com/pkg/errors"
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
//var Version = "0.0.0"

// Helloworld is the basic form of a gomake task.
// All gomake tasks return a closure to work with the "deps" API.
func Helloworld() func() {
	return func() {
		fmt.Println("Hello World")
	}
}

// Sub is an example parent command
func Sub() func() {
	return func() {
		fmt.Println("Sub says hi")
	}
}

// SubCmd is an example child command
func SubCmd() func() {
	return func() {
		fmt.Println("sub cmd says hi")
	}
}

// Hyphenated_cmd is an example of a hyphenated-cmd.
//
// Use snake case if your command has many words but should not be considered a
// child command.
func Hyphenated_cmd() func() {
	return func() {
		fmt.Println("Hyphenated_cmd says hi")
	}
}

// Hyphenated_cmdFoo_bar is an example of a child command of a hyphenated-cmd.
//
// Feel free to mix and match both hyphenated and non hyphenated command names.
// Yeah I know it's a pretty ugly function name but I don't expect this to be
// used as much as sub commands.
func Hyphenated_cmdFoo_bar() func() {
	return func() {
		fmt.Println("Hyphenated_cmdFoo_bar says hi")
	}
}

// Valid_params shows all valid parameter types
func Valid_params(
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
) func() {
	return func() {}
}

// Documented_flag is an example of how to add descriptions against a parameter.
//
// --foo: Is an example parameter
// Flag descriptions must only span a single line, this is not part of
// the --foo description.
func Documented_flag(foo string) func() {
	return func() {
		fmt.Println("Documentedflag = " + foo)
	}
}

// Short_flag is an example of how to set a shortned alias against a parameter.
//
// --foo, -f: Is an example parameter, with a short alias
func Short_flag(foo string) func() {
	return func() {
		fmt.Println("Shortflag = " + foo)
	}
}

// ErrorExample shows what happens when a task returns an error
//
// This is an example of returning a stdlib error. The fundamental problem with
// golang errors is that context is often lost. It is suggested to look at using
// one of:
//
//     - <https://github.com/pkg/errors>
//     - <https://github.com/go-errors/errors>
//
// gomake supports (through the use of <https://github.com/brad-jones/goerr>)
// returning stack traces from both these types of errors,
// see ErrorExample1 and ErrorExample2
//
// For more reading about this issue:
// <https://dave.cheney.net/tag/error-handling>
//
// > TIP: Set "GOMAKE_DEBUG=1" to see a stack trace.
func ErrorExample() func() error {
	return func() error {
		return errors.New("ops we failed for some reason")
	}
}

// ErrorExample1 is an example of using <https://github.com/pkg/errors>
//
// > TIP: Set "GOMAKE_DEBUG=1" to see a stack trace.
func ErrorExample1() func() error {
	return func() error {
		return errors1.New("ops we failed for some reason")
	}
}

// ErrorExample2 is an example of using <https://github.com/go-errors/errors>
//
// > TIP: Set "GOMAKE_DEBUG=1" to see a stack trace.
func ErrorExample2() func() error {
	return func() error {
		return errors2.New("ops we failed for some reason")
	}
}

// ErrorExample3 shows what happens when a task panics
//
// > TIP: Set "GOMAKE_DEBUG=1" to see a stack trace.
func ErrorExample3() func() {
	return func() {
		panic("ops we failed for some reason")
	}
}

// Go_async_example is the third type of function signature that is supported.
//
// gomake is built on top of the goasync library, see:
// <https://godoc.org/github.com/brad-jones/goasync>
//
// Use this style if you want more control over the execution of your tasks.
//
// > NOTE: We return `func(t *task.Internal)` and not `*task.Task` so that
// >       execution is deferred. This allows one to still execute these
// >       gomake tasks serially if desired.
func Go_async_example() func(t *task.Internal) {
	return func(t *task.Internal) {
		for i := 1; i < 5; i++ {
			if t.ShouldStop() {
				t.Reject(errors.New("got told to stop"))
				return
			}
			fmt.Println("running for the", i, "time")
			time.Sleep(1 * time.Second)
		}
		t.Resolve("finished work")
	}
}

// DepsExample shows how you can call other gomake tasks.
//
// The deps API is a wrapper around:
// <https://godoc.org/github.com/brad-jones/goasync/await>
//
// But adds in additional logic to ensure tasks are only executed once.
// For example if `a()`, depends on `b()` and `c()` depends on `b()` and your
// function depends on `a()` and `c()`, then `b()` will only execute once.
//
// The deps API will always fail fast, if you need more control over how your
// tasks are executed I suggest using the `await` library directly.
func DepsExample() func(t *task.Internal) {
	return deps.Serial(
		Sub(),
		SubCmd(),
		deps.Parallel(
			Short_flag(""),
			Documented_flag(""),
		),
	)
}
