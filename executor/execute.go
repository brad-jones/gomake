package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"gopkg.in/brad-jones/gomake.v2/generator"
)

func Execute(dir string, args ...string) error {

	dir, err := findGomakeFolder(dir)
	if err != nil {
		return err
	}

	reBuild := true

	exePath := filepath.Join(dir, "runner")
	if runtime.GOOS == "windows" {
		exePath = exePath + ".exe"
	}

	if _, err := os.Stat(exePath); err == nil {
		if hash1, err := generator.CacheHashRead(dir); err == nil {
			if hash2, err := generator.CacheHashGen(dir); err == nil {
				reBuild = hash1 != hash2
			}
		}
	}

	if reBuild {
		if err := generator.Generate(dir); err != nil {
			return err
		}
		cmd := exec.Command("go", "build", "-o", exePath, dir)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}
			return err
		}
	}

	// Ensure that the current working directory is one level above the
	// ".gomake" folder. This makes it easier to write tasks as you can
	// always rely on the working directory being the solution root for
	// any given project.
	//
	// NOTE: This obviously only changes the working dir for this process
	//       and has no effect on the users shell.
	if err := os.Chdir(filepath.Dir(dir)); err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		cmd := exec.Command(exePath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}
		}
		return err
	}

	newArgs := []string{exePath}
	newArgs = append(newArgs, args...)
	return syscall.Exec(exePath, newArgs, os.Environ())
}
