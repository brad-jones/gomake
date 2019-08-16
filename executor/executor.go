package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/brad-jones/goerr"
	"github.com/brad-jones/gomake/v3/generator"
)

func Execute(cwd string, args ...string) (err error) {
	defer goerr.Handle(func(e error) {
		if e != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}
		}
		err = e
	})

	// Find the folder that contains a valid gomake task runner
	dir := mustFindGomakeFolder(cwd)

	// Generate a hash of this task runner
	hash := generator.MustCacheHashGen(dir)

	// Create the path to the task runner executable
	homeDir, err := os.UserHomeDir()
	goerr.Check(err)
	exePath := filepath.Join(homeDir, ".gomake", hash)
	if runtime.GOOS == "windows" {
		exePath = exePath + ".exe"
	}

	// If it does not exist we need to build it
	if !fileExists(exePath) {
		generator.MustGenerate(dir)
		cmd := exec.Command("go", "build", "-o", exePath, filepath.Join(dir, ".gomake"))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		goerr.Check(cmd.Run())
		goerr.Check(os.RemoveAll(filepath.Join(dir, ".gomake")))
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
		goerr.Check(cmd.Run())
	} else {
		newArgs := []string{exePath}
		newArgs = append(newArgs, args...)
		goerr.Check(syscall.Exec(exePath, newArgs, os.Environ()))
	}

	return err
}

// ErrReachedRootOfFs is returned when findGoMakeFolder can not find a
// valid "makefile.go" file after having recursed all the way up to the
// root of filesystem.
type ErrReachedRootOfFs struct {
}

func (e *ErrReachedRootOfFs) Error() string {
	return "gomake: failed to find a valid 'makefile.go', reached root of filesystem"
}

func findGomakeFolder(dir string) (string, error) {
	goMakeFile := filepath.Join(dir, "makefile.go")
	if _, err := os.Stat(goMakeFile); err == nil {
		return dir, nil
	}
	if dir == filepath.VolumeName(dir)+string(os.PathSeparator) {
		return "", &ErrReachedRootOfFs{}
	}
	parentDir := filepath.Join(dir, "..")
	return findGomakeFolder(parentDir)
}

func mustFindGomakeFolder(dir string) string {
	v, e := findGomakeFolder(dir)
	goerr.Check(e)
	return v
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
