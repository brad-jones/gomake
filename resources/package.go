package resources

import (
	"io"
	"io/ioutil"

	"github.com/phogolabs/parcello"
)

//go:generate parcello -r

// TODO: You could then have a second "go:generate"
// command that generated a set of strongly typed
// functions like these... feel like another golang package comming along

// TypedFile provides 3 easy helper functions to
// easily read an embedded parcello resource.
// It also returns the original file path.
type TypedFile struct {
	Reader func() io.Reader
	Bytes  func() []byte
	String func() string
	Path   string
}

func typedFileFactory(path string) *TypedFile {
	file := &TypedFile{Path: path}
	file.Reader = func() io.Reader {
		file, err := parcello.Open(path)
		if err != nil {
			panic(err)
		}
		return file
	}
	file.Bytes = func() []byte {
		dat, err := ioutil.ReadAll(file.Reader())
		if err != nil {
			panic(err)
		}
		return dat
	}
	file.String = func() string {
		return string(file.Bytes())
	}
	return file
}
