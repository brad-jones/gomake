package generator

import (
	"crypto/sha1"
	"fmt"
	"go/token"
	"os"
	"strconv"
	"strings"
)

func CacheHashGen(dir string) (string, error) {
	fset, _, err := parseAST(dir)
	if err != nil {
		return "", err
	}
	h, err := cacheHashGen(fset)
	if err != nil {
		return "", err
	}
	return h, nil
}

func cacheHashGen(fset *token.FileSet) (string, error) {

	var possibleErr error
	var hashContent strings.Builder

	fset.Iterate(func(f *token.File) bool {
		fname := f.Name()
		finfo, err := os.Stat(fname)
		if err != nil {
			possibleErr = err
			return false
		}
		hashContent.WriteString(fname)
		hashContent.WriteString(strconv.FormatInt(finfo.ModTime().Unix(), 10))
		return true
	})

	if possibleErr != nil {
		return "", possibleErr
	}

	h := sha1.New()
	_, err := h.Write([]byte(hashContent.String()))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
