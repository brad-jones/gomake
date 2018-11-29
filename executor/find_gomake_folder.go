package executor

import (
	"os"
	"path/filepath"
)

func findGomakeFolder(dir string) (string, error) {
	goMakeFolder := filepath.Join(dir, ".gomake")
	if _, err := os.Stat(goMakeFolder); err == nil {
		return goMakeFolder, nil
	}
	if dir == filepath.VolumeName(dir)+string(os.PathSeparator) {
		return "", &ErrReachedRootOfFs{}
	}
	parentDir := filepath.Join(dir, "..")
	return findGomakeFolder(parentDir)
}
