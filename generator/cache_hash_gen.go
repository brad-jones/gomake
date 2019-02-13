package generator

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func CacheHashGen(dir string) (string, error) {

	var hashContent strings.Builder

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") && info.Name() != "makefile_generated.go" {
			finfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			hashContent.WriteString(path)
			hashContent.WriteString(strconv.FormatInt(finfo.ModTime().Unix(), 10))
		}
		return nil
	}); err != nil {
		return "", err
	}

	h := sha1.New()
	_, err := h.Write([]byte(hashContent.String()))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
