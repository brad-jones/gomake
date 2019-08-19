package executor

// Generator is something that can generate the complete source code for a
// CLI application based on the exported functions from the package of a
// "makefile.go" file.
type Generator interface {
	// MustGenApp should accept a path to a folder where a valid "makefile.go"
	// is found and then return a path of the generated source code.
	//
	// eg: ".gomake/main.go", where ".gomake" is a new folder contained in
	// the same parent as "makefile.go". This folder is what will be built
	// into a binary and then deleted.
	MustGenApp(path string) string

	// MustGenAppHash should accept a path to a folder where a valid "makefile.go"
	// is found and then create a hash that represents the entire source code that
	// makes up the package exported from "makefile.go".
	//
	// This hash is used to determin if the task runner needs to be re-complied or not.
	MustGenAppHash(path string) string
}
