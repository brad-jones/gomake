package generator

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

func CacheHashRead(dir string) (string, error) {

	file, err := os.Open(filepath.Join(dir, "makefile_generated.go"))
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	scanner.Scan()
	line := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	hash := strings.TrimSpace(strings.Replace(line, "// cache-hash: ", "", 1))

	return hash, nil
}
