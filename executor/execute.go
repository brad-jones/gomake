package executor

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/brad-jones/gomake/generator"
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

	if runtime.GOOS == "windows" {
		cmd := exec.Command(exePath, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Exit(1)
			}
			return err
		}
	}

	newArgs := []string{exePath}
	newArgs = append(newArgs, args...)
	return syscall.Exec(exePath, newArgs, os.Environ())
}
