package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

func findGomakeFolder(dir string) (string, error) {
	goMakeFolder := filepath.Join(dir, ".gomake")
	if _, err := os.Stat(goMakeFolder); err == nil {
		return goMakeFolder, nil
	}
	// TODO: Consider root detection for Windows
	fmt.Println("DIR: " + dir)
	if dir == "/" {
		return "", &ErrReachedRootOfFs{}
	}
	parentDir := filepath.Join(dir, "..")
	return findGomakeFolder(parentDir)
}
