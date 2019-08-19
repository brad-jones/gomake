package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/brad-jones/goerr"
	"github.com/go-errors/errors"
)

// Executor is a class like object, create new instances with `executor.New()`
type Executor struct {
	searchDir string
	g         Generator
}

// New creates an instance of `Executor`
func New(g Generator) *Executor {
	return &Executor{
		g: g,
	}
}

// SetSearchDir accepts a path that will be used as the starting point for
// looking for a "makefile.go". If not set this will default to the current
// working directory.
func (e *Executor) SetSearchDir(path string) *Executor {
	e.searchDir = path
	return e
}

// GetSearchDir a path that will be used as the starting point for looking for
// a "makefile.go". If not set this will default to the current working directory.
func (e *Executor) GetSearchDir() (path string, err error) {
	defer goerr.Handle(func(e error) {
		path = ""
		err = e
	})
	if e.searchDir == "" {
		dir, err := os.Getwd()
		goerr.Check(err)
		e.searchDir = dir
	}
	path = e.searchDir
	return
}

// MustGetSearchDir does the same thing as GetSearchDir but panics if an error is encountered
func (e *Executor) MustGetSearchDir() string {
	path, err := e.GetSearchDir()
	goerr.Check(err)
	return path
}

// Execute orchestrates the process of generating, compiling and running a
// gomake "makefile.go" task runner. It takes the arguments that will be passed
// on to the compiled task runner.
func (e *Executor) Execute(args ...string) (err error) {
	defer goerr.Handle(func(e error) {
		err = errors.Wrap(e, 0)
	})

	dir := e.mustFindGomakeFolder(e.MustGetSearchDir())
	hash := e.g.MustGenAppHash(dir)
	exePath := e.mustGetTaskRunnerExePath(hash)

	// We need to build a new task runner to execute
	if !fileExists(exePath) {
		appPath := e.g.MustGenApp(dir)
		cmd := exec.Command("go", "build", "-o", exePath, appPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		goerr.Check(cmd.Run())
		if os.Getenv("GOMAKE_PARSE_DEBUG") != "1" {
			goerr.Check(os.RemoveAll(appPath))
		}
	}

	// Ensure that the current working directory is the same as the
	// "makefile.go" file. This makes it easier to write tasks as you can
	// always rely on the working directory being the solution root for
	// any given project.
	//
	// NOTE: This obviously only changes the working dir for this process
	//       and has no effect on the users shell.
	goerr.Check(os.Chdir(dir))

	// Now lets execute the task runner executable, for windows we just stream
	// STDIN, STDOUT and STDERR. For unix systems we can pass execution directly
	// to the task runner binary.
	if runtime.GOOS == "windows" {
		cmd := exec.Command(exePath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}
			goerr.Check(err)
		}
		os.Exit(0)
	}

	newArgs := []string{exePath}
	newArgs = append(newArgs, args...)
	goerr.Check(syscall.Exec(exePath, newArgs, os.Environ()))

	return
}

// MustExecute does the same thing as Execute but panics if an error is encountered
func (e *Executor) MustExecute(args ...string) {
	goerr.Check(e.Execute(args...))
}

func (e *Executor) findGomakeFolder(dir string) (string, error) {
	goMakeFile := filepath.Join(dir, "makefile.go")
	if _, err := os.Stat(goMakeFile); err == nil {
		return dir, nil
	}
	if dir == filepath.VolumeName(dir)+string(os.PathSeparator) {
		return "", errors.New(ErrReachedRootOfFs)
	}
	parentDir := filepath.Join(dir, "..")
	return e.findGomakeFolder(parentDir)
}

func (e *Executor) mustFindGomakeFolder(dir string) string {
	v, err := e.findGomakeFolder(dir)
	goerr.Check(err)
	return v
}

func (e *Executor) getTaskRunnerExePath(hash string) (path string, err error) {
	defer goerr.Handle(func(e error) {
		path = ""
		err = e
	})

	homeDir, err := os.UserHomeDir()
	goerr.Check(err)

	path = filepath.Join(homeDir, ".gomake", hash)
	if runtime.GOOS == "windows" {
		path = path + ".exe"
	}

	return
}

func (e *Executor) mustGetTaskRunnerExePath(hash string) string {
	v, err := e.getTaskRunnerExePath(hash)
	goerr.Check(err)
	return v
}
