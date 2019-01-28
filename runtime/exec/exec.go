package exec

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
)

// Cmd provides a fluent, decorator based API for os/exec.
// Eg: Cmd("ping", Args("-c", "4", "1.1.1.1"))
func Cmd(cmd string, decorators ...func(*exec.Cmd) error) (*exec.Cmd, error) {
	c := exec.Command(cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	for _, decorator := range decorators {
		if err := decorator(c); err != nil {
			return nil, err
		}
	}

	if c.Env == nil {
		c.Env = os.Environ()
	}

	return c, nil
}

// Run is a convenience function for simple cases.
// Instead of: Cmd("ping", Args("-c", "4", "8.8.8.8")).Run()
// You might write: Run("ping", "-c", "4", "8.8.8.8")
//
// NOTE: A closure is returned thus making for nicer syntax when used with
// run.Serial/Parallel functions. If using standalone you will need to invoke
// the returned closure.
func Run(cmd string, args ...string) func() error {
	return func() error {
		c, err := Cmd(cmd, Args(args...))
		if err != nil {
			return err
		}
		return c.Run()
	}
}

// RunBuffered is a convenience function for simple cases.
// Instead of: RunBufferedCmd(Cmd("ping", Args("-c", "4", "8.8.8.8")))
// You might write: RunBuffered("ping", "-c", "4", "8.8.8.8")
func RunBuffered(cmd string, args ...string) (stdOutBuf, stdErrBuf string, err error) {
	c, err := Cmd(cmd, Args(args...))
	if err != nil {
		return "", "", err
	}
	o, e, err := RunBufferedCmd(c)
	return string(o), string(e), err
}

// RunPrefixed is a convenience function for simple cases.
// Instead of: RunPrefixedCmd("foo", Cmd("ping", Args("-c", "4", "8.8.8.8")))
// You might write: RunPrefixed("foo", "ping", "-c", "4", "8.8.8.8")
//
// NOTE: A closure is returned thus making for nicer syntax when used with
// run.Serial/Parallel functions. If using standalone you will need to invoke
// the returned closure.
func RunPrefixed(prefix, cmd string, args ...string) func() error {
	return func() error {
		c, err := Cmd(cmd, Args(args...))
		if err != nil {
			return err
		}
		return RunPrefixedCmd(prefix, c)
	}
}

// RunPrefixedCmd will prefix all StdOut and StdErr with given prefix.
// This is useful when running many commands concurrently,
// output will look similar to docker-compose.
func RunPrefixedCmd(prefix string, cmd *exec.Cmd) error {
	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	if err != nil {
		return err
	}
	cmd.Stdout = stdOutPipeW

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	if err != nil {
		return err
	}
	cmd.Stderr = stdErrPipeW

	return prefixed(prefix,
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return cmd.Run()
		},
	)
}

// RunBufferedCmd will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a command.
func RunBufferedCmd(cmd *exec.Cmd) (stdOutBuf, stdErrBuf []byte, err error) {
	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	cmd.Stdout = stdOutPipeW

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	cmd.Stderr = stdErrPipeW

	return buffered(
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return cmd.Run()
		},
	)
}

// Pipe will send the output of the first command
// to the input of the second and so on.
func Pipe(cmds ...*exec.Cmd) error {
	for key, cmd := range cmds {
		if key > 0 {
			cmds[key-1].Stdout = nil
			if pipe, err := cmds[key-1].StdoutPipe(); err == nil {
				cmd.Stdin = pipe
			} else {
				return err
			}
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return err
		}
	}

	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return err
		}
	}

	return nil
}

// PipePrefixed will prefix all StdOut and StdErr with given prefix.
// This is useful when running many pipes concurrently,
// output will look similar to docker-compose.
func PipePrefixed(prefix string, cmds ...*exec.Cmd) error {
	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	if err != nil {
		return err
	}

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	if err != nil {
		return err
	}

	for _, cmd := range cmds {
		cmd.Stderr = stdErrPipeW
	}

	cmds[len(cmds)-1].Stdout = stdOutPipeW

	return prefixed(prefix,
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return Pipe(cmds...)
		},
	)
}

// PipeBuffered will buffer all StdOut and StdErr, returning the buffers.
// This is useful when you wish to parse the results of a pipe.
func PipeBuffered(cmds ...*exec.Cmd) (stdOutBuf, stdErrBuf []byte, err error) {
	stdOutPipeR, stdOutPipeW, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	stdErrPipeR, stdErrPipeW, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	for _, cmd := range cmds {
		cmd.Stderr = stdErrPipeW
	}

	cmds[len(cmds)-1].Stdout = stdOutPipeW

	return buffered(
		os.Stdout, os.Stderr,
		stdOutPipeR, stdErrPipeR,
		stdOutPipeW, stdErrPipeW, func() error {
			return Pipe(cmds...)
		},
	)
}

// Args allows you to define the arguments sent to the command to be run
func Args(args ...string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Args = append([]string{c.Path}, args...)
		return nil
	}
}

// Cwd allows you to configure the working directory of the command to be run
func Cwd(dir string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Dir = dir
		return nil
	}
}

// Env allows you to set the exact environment in which the command will be run
func Env(env map[string]string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		e := []string{}
		for k, v := range env {
			e = append(e, k+"="+v)
		}
		c.Env = e
		return nil
	}
}

// EnvCombined will add the variables you provide to the existing environment
func EnvCombined(env map[string]string) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		e := os.Environ()
		for k, v := range env {
			e = append(e, k+"="+v)
		}
		c.Env = e
		return nil
	}
}

// SetIn allows you to set a custom StdIn stream
func SetIn(in io.Reader) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stdin = in
		return nil
	}
}

// SetOut allows you to set a custom StdOut stream
func SetOut(out io.Writer) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stdout = out
		return nil
	}
}

// SetErr allows you to set a custom StdErr stream
func SetErr(err io.Writer) func(*exec.Cmd) error {
	return func(c *exec.Cmd) error {
		c.Stderr = err
		return nil
	}
}

func buffered(fstdOut, stdErr io.Writer,
	stdOutPipeR, stdErrPipeR io.Reader,
	stdOutPipeW, stdErrPipeW io.WriteCloser,
	fn func() error) (stdOutBuf, stdErrBuf []byte, err error) {

	errorCh := make(chan error)

	// Run the function, pipeing all StdOut and StdErr to our buffers
	go func() {
		defer stdOutPipeW.Close()
		defer stdErrPipeW.Close()
		errorCh <- fn()
	}()

	// Read all StdOut into our buffer
	stdOutC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdOutPipeR); err != nil {
			errorCh <- err
		} else {
			stdOutC <- buf.Bytes()
		}
	}()

	// Read all StdErr into our buffer
	stdErrC := make(chan []byte)
	go func() {
		var buf bytes.Buffer
		if _, err := io.Copy(&buf, stdErrPipeR); err != nil {
			errorCh <- err
		} else {
			stdErrC <- buf.Bytes()
		}
	}()

	// Catch any errors
	if err := <-errorCh; err != nil {
		return <-stdOutC, <-stdErrC, err
	}

	return <-stdOutC, <-stdErrC, nil
}

var prefixToColorMap = &sync.Map{}
var choosenColors = &sync.Map{}
var choosenCount = 0
var availableColors = []color.Attribute{
	color.FgRed,
	color.FgGreen,
	color.FgYellow,
	color.FgBlue,
	color.FgMagenta,
	color.FgCyan,
	color.FgHiRed,
	color.FgHiGreen,
	color.FgHiYellow,
	color.FgHiBlue,
	color.FgHiMagenta,
	color.FgHiCyan,
	color.BgRed,
	color.BgGreen,
	color.BgYellow,
	color.BgBlue,
	color.BgMagenta,
	color.BgCyan,
	color.BgHiRed,
	color.BgHiGreen,
	color.BgHiYellow,
	color.BgHiBlue,
	color.BgHiMagenta,
	color.BgHiCyan,
}

var randGen *rand.Rand

func init() {
	randGen = rand.New(rand.NewSource(time.Now().Unix()))
}

func colorChooser(prefix string) color.Attribute {

	// Return the cached value if it exists
	c, exists := prefixToColorMap.Load(prefix)
	if exists {
		return c.(color.Attribute)
	}

	// Randomly choose a new color
	c = availableColors[randGen.Intn(len(availableColors))]

	// Check if we reached the maximum number of available colors.
	// If so we will just have to reuse a color.
	if choosenCount == len(availableColors) {
		prefixToColorMap.Store(prefix, c)
		return c.(color.Attribute)
	}

	// Check if this color has already been choosen,
	// running ourselves again to select hopefully a different color.
	if _, choosen := choosenColors.Load(c); choosen == true {
		return colorChooser(prefix)
	}

	// Cache the result for next time
	choosenCount = choosenCount + 1
	choosenColors.Store(c, true)
	prefixToColorMap.Store(prefix, c)

	return c.(color.Attribute)
}

func prefixed(prefix string,
	stdOut, stdErr io.Writer,
	stdOutPipeR, stdErrPipeR io.Reader,
	stdOutPipeW, stdErrPipeW io.WriteCloser,
	fn func() error) error {

	errorCh := make(chan error)
	stdOutScanner := bufio.NewScanner(stdOutPipeR)
	stdErrScanner := bufio.NewScanner(stdErrPipeR)

	// Make the prefix colorful
	prefix = color.New(colorChooser(prefix)).Sprint(prefix + " | ")

	// Run the function, pipeing all StdOut and StdErr to our scanners
	go func() {
		defer stdOutPipeW.Close()
		defer stdErrPipeW.Close()
		errorCh <- fn()
	}()

	// Prefix all StdOut
	go func() {
		for stdOutScanner.Scan() {
			fmt.Fprintln(stdOut, prefix+strings.TrimSpace(stdOutScanner.Text())+"\r")
		}
		if err := stdOutScanner.Err(); err != nil {
			errorCh <- fmt.Errorf("prefxing standard out: %s", err)
		}
	}()

	// Prefix all StdErr
	go func() {
		for stdErrScanner.Scan() {
			fmt.Fprintln(stdErr, prefix+color.RedString(strings.TrimSpace(stdErrScanner.Text()))+"\r")
		}
		if err := stdErrScanner.Err(); err != nil {
			errorCh <- fmt.Errorf("prefxing standard err: %s", err)
		}
	}()

	// Catch any errors
	if err := <-errorCh; err != nil {
		return err
	}

	return nil
}
