package generator

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/brad-jones/goerr"
	"github.com/brad-jones/gomake/v3/resources"
	"github.com/go-errors/errors"
	"github.com/pinzolo/casee"
	"github.com/sirkon/goproxy/gomod"
	git "gopkg.in/src-d/go-git.v4"
)

// Generator is a class like object, create new instances with `generator.New()`
type Generator struct {
	fset *token.FileSet
}

// New creates an instance of `Generator`
func New() *Generator {
	return &Generator{
		fset: token.NewFileSet(),
	}
}

// MustGenApp should accept a path to a folder where a valid "makefile.go"
// is found and then return a path of the generated source code.
//
// eg: ".gomake/main.go", where ".gomake" is a new folder contained in
// the same parent as "makefile.go". This folder is what will be built
// into a binary and then deleted.
func (g *Generator) MustGenApp(path string) string {
	packageName, nodes := g.mustParseAST(path)
	viewModel := g.mustBuildViewModel(packageName, nodes)
	return g.mustWriteGeneratedMakefile(path, viewModel)
}

// MustGenAppHash should accept a path to a folder where a valid "makefile.go"
// is found and then create a hash that represents the entire source code that
// makes up the package exported from "makefile.go".
//
// This hash is used to determin if the task runner needs to be re-complied or not.
func (g *Generator) MustGenAppHash(path string) string {
	var hashContent strings.Builder

	// First of all we take into account the go.mod (& go.sum if it exists)
	// This means the moment you import a new package into your task runner
	// or update an existing one the task runner will get rebuilt.
	goModPath, goSumPath := g.mustFindGoModSumFiles(path)
	goModInfo, err := os.Stat(goModPath)
	goerr.Check(err)
	hashContent.WriteString(goModPath)
	hashContent.WriteString(strconv.FormatInt(goModInfo.ModTime().Unix(), 10))
	if goSumPath != "" {
		goSumInfo, err := os.Stat(goSumPath)
		goerr.Check(err)
		hashContent.WriteString(goSumPath)
		hashContent.WriteString(strconv.FormatInt(goSumInfo.ModTime().Unix(), 10))
	}

	// Now we take into account each source file "local" to the go module.
	// So regardless of how you structure your task runner any changes will be detected.
	goMod := g.mustParseModFile(goModPath)
	var recurseAst func(path string)
	recurseAst = func(path string) {
		_, nodes := g.mustParseAST(path)
		for srcFilePath, fileAst := range nodes.Files {
			fInfo, err := os.Stat(srcFilePath)
			goerr.Check(err)
			hashContent.WriteString(srcFilePath)
			hashContent.WriteString(strconv.FormatInt(fInfo.ModTime().Unix(), 10))
			for _, importDec := range fileAst.Imports {
				importPath := strings.ReplaceAll(importDec.Path.Value, "\"", "")
				if strings.HasPrefix(importPath, goMod.Name) {
					recurseAst(filepath.Join(filepath.Dir(goModPath), strings.Replace(importPath, goMod.Name, "", 1)))
				}
			}
		}
	}
	recurseAst(path)

	// Generate and return a SHA1 hash
	h := sha1.New()
	_, err = h.Write([]byte(hashContent.String()))
	goerr.Check(err)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (g *Generator) mustParseAST(cwd string) (packageName string, parsedAST *ast.Package) {
	ignoreGeneratedMakefile := func(file os.FileInfo) bool {
		return file.Name() != "makefile_generated.go"
	}

	pkgs, err := parser.ParseDir(g.fset, cwd, ignoreGeneratedMakefile, parser.ParseComments)
	goerr.Check(err)

	for k, v := range pkgs {
		if !strings.HasSuffix(k, "_test") {
			if os.Getenv("GOMAKE_PARSE_DEBUG") == "1" {
				ast.Print(g.fset, v)
			}
			return k, v
		}
	}

	goerr.Check(errors.New(ErrNoPackageNotFound))
	return
}

// Command names can not start with a number.
// What we are trying to avoid is command invocations like:
//
//     gomake upload-to-s 3
//
// It should obviously be:
//
//     gomake upload-to-s3
var cmdStartsWithNo = regexp.MustCompile(".*?(_\\d+).*?")

func (g *Generator) mustBuildCmdTree(nodes *ast.Package) cmdTree {

	tree := cmdTree{}

	ast.Inspect(nodes, func(node ast.Node) bool {
		switch funcDecl := node.(type) {
		case *ast.FuncDecl:
			if g.isValidGoMakeFunc(funcDecl) {
				root := tree
				cmdName := strings.Replace(funcDecl.Name.Name, "_", ":", -1)
				cmdName = casee.ToSnakeCase(cmdName)
				cmdName = strings.Replace(cmdName, "_:_", "-", -1)
				for _, match := range cmdStartsWithNo.FindAllStringSubmatch(cmdName, -1) {
					cmdName = strings.Replace(cmdName, match[1], strings.TrimPrefix(match[1], "_"), 1)
				}
				cmdParts := strings.Split(cmdName, "_")
				cmdPartsLen := len(cmdParts)
				for i, cmd := range cmdParts {
					if _, ok := root[cmd]; !ok {
						if i+1 == cmdPartsLen {
							root[cmd] = cmdTreeBranch{
								leaf:   funcDecl,
								branch: cmdTree{},
							}
						} else {
							root[cmd] = cmdTreeBranch{
								leaf: &ast.FuncDecl{
									Name: &ast.Ident{
										Name: "gomake_noop",
									},
									Type: &ast.FuncType{},
								},
								branch: cmdTree{},
							}
						}
					}
					root = root[cmd].branch
				}
			}
		}
		return true
	})

	return tree
}

func (g *Generator) findGoModSumFiles(dir string) (goModPath, goSumPath string, err error) {
	goModPath = filepath.Join(dir, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		goSumPath = filepath.Join(dir, "go.sum")
		if _, err := os.Stat(goSumPath); err == nil {
			return goModPath, goSumPath, nil
		}
		return goModPath, "", nil
	}
	if dir == filepath.VolumeName(dir)+string(os.PathSeparator) {
		return "", "", errors.New(ErrReachedRootOfFs)
	}
	parentDir := filepath.Join(dir, "..")
	return g.findGoModSumFiles(parentDir)
}

func (g *Generator) mustFindGoModSumFiles(dir string) (goModPath, goSumPath string) {
	goModPath, goSumPath, err := g.findGoModSumFiles(dir)
	goerr.Check(err)
	return
}

func (g *Generator) mustParseModFile(path string) *gomod.Module {
	dat, err := ioutil.ReadFile(path)
	goerr.Check(err)
	mod, err := gomod.Parse(path, dat)
	goerr.Check(err)
	return mod
}

func (g *Generator) mustGetModulePath(nodes *ast.Package) string {
	for srcFilePath := range nodes.Files {
		goModPath, _ := g.mustFindGoModSumFiles(filepath.Dir(srcFilePath))
		mod := g.mustParseModFile(goModPath)
		goModDir := filepath.Dir(goModPath)
		srcFileDir := filepath.Dir(srcFilePath)
		diff := strings.Replace(srcFileDir, goModDir, "", 1)
		return mod.Name + diff
	}
	panic(errors.New("gomake: could not get go.mod path"))
}

func (g *Generator) mustBuildViewModel(packageName string, nodes *ast.Package) *tplViewModel {
	viewModel := &tplViewModel{
		MakefilePackageImport: packageName + " \"" + g.mustGetModulePath(nodes) + "\"",
	}
	cmdTree := g.mustBuildCmdTree(nodes)
	viewModel.Commands = g.mustWalkCmdTree(cmdTree, packageName, "root", "", 0)

	ast.Inspect(nodes, func(node ast.Node) bool {
		switch n := node.(type) {
		case *ast.ValueSpec:
			switch n.Names[0].Name {
			case "Use":
				viewModel.Use = strings.Trim(n.Values[0].(*ast.BasicLit).Value, "\"")
			case "Short":
				viewModel.Short = strings.Trim(n.Values[0].(*ast.BasicLit).Value, "\"")
			case "Version":
				viewModel.Version = strings.Trim(n.Values[0].(*ast.BasicLit).Value, "\"")
			}
		}
		return true
	})

	if viewModel.Use == "" {
		viewModel.Use = "gomake"
	}

	if viewModel.Short == "" {
		viewModel.Short = "Makefile written in golang"
	}

	if viewModel.Version == "" {
		viewModel.Version = "0.0.0"
		repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
		if err == nil {
			ref, err := repo.Head()
			if err == nil {
				viewModel.Version = "0.0.0-" + ref.Name().Short() + "-" + ref.Hash().String()
			}
		}
	}

	return viewModel
}

var regexpOptDoc = regexp.MustCompile("--(.*?),?:")

func (g *Generator) mustWalkCmdTree(tree cmdTree, packageName, parentCmdName, fullCmdName string, depth int) []*tplCommand {

	cmds := []*tplCommand{}

	for cmdName, branch := range tree {

		// Add a new command
		cmd := &tplCommand{}
		cmd.CmdName = cmdName
		cmd.FullCmdName = strings.TrimSpace(fullCmdName + " " + cmdName)
		cmd.CmdDepth = depth

		if branch.leaf.Name.Name != "gomake_noop" {
			cmd.FuncName = packageName + "." + branch.leaf.Name.Name
		} else {
			cmd.FuncName = branch.leaf.Name.Name
		}

		cmd.CobraCmdName = casee.ToCamelCase(cmdName)
		cmd.ParentCmdName = casee.ToCamelCase(parentCmdName + "Cmd")
		cmds = append(cmds, cmd)

		// Add any sub commands
		cmd.Commands = g.mustWalkCmdTree(branch.branch, packageName, cmdName, cmd.FullCmdName, depth+1)

		// If the function has a doc comment we will use it to provide
		// the long and short descriptions of the cobra command.
		// The first line being the short description, with anything
		// following making up the remainder of the long description.
		optDocs := map[string]string{}
		optShortNames := map[string]string{}
		optDefaultValues := map[string]string{}
		if branch.leaf.Doc != nil {
			for k, v := range branch.leaf.Doc.List {
				line := strings.TrimLeft(v.Text, "// ")
				line = strings.Replace(line, "\"", "\\\"", -1)
				if k == 0 {
					cmd.ShortDescription = line
				} else {
					// Lines that start with "--" document command
					// options they should not be included in the
					// overall command documentation
					if !strings.HasPrefix(line, "--") {
						cmd.LongDescription = cmd.LongDescription + "\\n" + line
					} else {
						// Parse option comments, they look like this:
						// --foo, -f: the foo option does xyz
						// --bar: the bar option does abc
						// --baz, $BAZ: yes defaults can be provided via the env
						// --abc, -a, $ABC: also supported
						// --qux, "a static default": this should also work
						// They may only span a single line
						parts := strings.Split(regexpOptDoc.FindStringSubmatch(line)[1], ",")
						value := strings.TrimSpace(regexpOptDoc.ReplaceAllString(line, ""))
						optDocs[parts[0]] = value
						if len(parts) > 1 {
							v := strings.TrimSpace(parts[1])
							if strings.HasPrefix(v, "-") {
								optShortNames[parts[0]] = strings.TrimPrefix(v, "-")
							} else {
								optDefaultValues[parts[0]] = v
							}
							if len(parts) > 2 {
								optDefaultValues[parts[0]] = strings.TrimSpace(parts[2])
							}
						}
					}
				}
			}
		}

		cmd.LongDescription = "\\n" + strings.TrimLeft(cmd.LongDescription, "\\n")
		if cmd.LongDescription == "\\n" {
			cmd.LongDescription = ""
		}

		// Function parameters are converted to cobra command options.
		// A final variadic parameter is converted to cobra positional arguments
		// NOTE: Only `...string` is supported
		cmd.Options = []*tplOption{}
		cmd.Args = &tplArgs{NoArgs: true}
		if branch.leaf.Type.Params != nil {
			for _, param := range branch.leaf.Type.Params.List {

				typeName, isArray, isPositional := g.mustGetTypeFromAstNode(param)

				if typeName == "Context" {
					cmd.HasCtx = true
					continue
				}

				if isPositional {
					if typeName != "string" {
						goerr.Check(errors.New(ErrVariadicMustBeString))
					}
					cmd.Args.NoArgs = false
					continue
				}

				flagType := g.mustMapTypeToFlagType(typeName, isArray)
				defaultValue := g.mustGetDefaultValueForType(typeName, isArray)

				for _, name := range param.Names {
					optNameChain := casee.ToChainCase(name.Name)
					description := optDocs[optNameChain]
					optDefaultValue := defaultValue
					if val, ok := optDefaultValues[optNameChain]; ok {
						if strings.HasPrefix(val, "$") {
							description = "env: " + val + " - " + description
							val = "os.Getenv(\"" + strings.TrimPrefix(val, "$") + "\")"
						}
						optDefaultValue = val
					}
					optDefaultValue = strings.Replace(optDefaultValue, "\\\"", "\"", -1)
					cmd.Options = append(cmd.Options, &tplOption{
						Name:         name.Name,
						LongName:     optNameChain,
						ShortName:    optShortNames[optNameChain],
						Description:  description,
						DefaultValue: optDefaultValue,
						FlagType:     flagType,
					})
				}
			}
		}

		if branch.leaf.Type.Results != nil {
			cmd.HasErr = true
		}
	}

	return cmds
}

func (g *Generator) mustGetDefaultValueForType(typeName string, isArray bool) string {
	if isArray {
		return "nil"
	}

	switch typeName {
	case "string":
		return `""`
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"Duration":
		return "0"
	case "bool":
		return "false"
	}

	return "nil"
}

func (g *Generator) mustGetTypeFromAstNode(node *ast.Field) (typeName string, isArray bool, isPositional bool) {
	switch n := node.Type.(type) {
	case *ast.Ident:
		typeName = n.Name

	case *ast.SelectorExpr:
		typeName = n.Sel.Name

	case *ast.ArrayType:
		switch et := n.Elt.(type) {
		case *ast.Ident:
			typeName = et.Name
		case *ast.SelectorExpr:
			typeName = et.Sel.Name
		}
		isArray = true

	case *ast.Ellipsis:
		switch et := n.Elt.(type) {
		case *ast.Ident:
			typeName = et.Name
		case *ast.SelectorExpr:
			typeName = et.Sel.Name
		}
		isArray = true
		isPositional = true
	}

	return typeName, isArray, isPositional
}

func (g *Generator) isValidGoMakeFunc(funcDecl *ast.FuncDecl) bool {
	firstFuncChar := string(funcDecl.Name.Name[0])
	if strings.ToUpper(firstFuncChar) == firstFuncChar {
		if funcDecl.Type.Results == nil {
			return true
		}
		if len(funcDecl.Type.Results.List) == 1 {
			if r, ok := funcDecl.Type.Results.List[0].Type.(*ast.Ident); ok {
				if r.Name == "error" {
					return true
				}
			}
		}
	}
	return false
}

var supportedFlagTypes = map[string]struct{}{
	"BoolP":          {},
	"BoolSliceP":     {},
	"StringP":        {},
	"StringSliceP":   {},
	"IntP":           {},
	"IntSliceP":      {},
	"Int8P":          {},
	"Int16P":         {},
	"Int32P":         {},
	"Int64P":         {},
	"Float32P":       {},
	"Float64P":       {},
	"UintP":          {},
	"UintSliceP":     {},
	"Uint8P":         {},
	"Uint16P":        {},
	"Uint32P":        {},
	"Uint64P":        {},
	"DurationP":      {},
	"DurationSliceP": {},
	"IPP":            {},
	"IPSliceP":       {},
	"IPMaskP":        {},
	"BytesBase64P":   {},
}

func (g *Generator) mustMapTypeToFlagType(typeName string, isArray bool) (flagType string) {

	if typeName == "IP" || typeName == "IPMask" {
		flagType = typeName
	} else {
		flagType = casee.ToPascalCase(typeName)
	}

	if isArray {
		flagType = flagType + "SliceP"
	} else {
		flagType = flagType + "P"
	}

	if flagType == "ByteSliceP" {
		flagType = "BytesBase64P"
	}

	if _, exists := supportedFlagTypes[flagType]; !exists {
		goerr.Check(errors.New(&ErrUnsupportedFlagType{
			OriginalTypeName: typeName,
			IsArray:          isArray,
			MappedTypeName:   flagType,
		}))
	}

	return flagType
}

func (g *Generator) mustWriteGeneratedMakefile(cwd string, viewModel *tplViewModel) string {
	tpl, err := template.New("generated").
		Parse(resources.Typed.MakefileGeneratedGotmpl.String())
	goerr.Check(err)

	var unFormatted bytes.Buffer
	goerr.Check(tpl.Execute(&unFormatted, viewModel))

	formatted, err := format.Source(unFormatted.Bytes())
	goerr.Check(err)

	appPath := filepath.Join(cwd, ".gomake")
	goerr.Check(os.MkdirAll(appPath, 0744))
	goerr.Check(ioutil.WriteFile(filepath.Join(appPath, "makefile_generated.go"), formatted, 0644))

	return appPath
}
