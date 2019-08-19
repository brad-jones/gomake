package generator

import "go/ast"

type cmdTree map[string]cmdTreeBranch

type cmdTreeBranch struct {
	leaf   *ast.FuncDecl
	branch cmdTree
}

type tplViewModel struct {
	Use                   string
	Short                 string
	Version               string
	MakefilePackageImport string
	Commands              []*tplCommand
}

type tplCommand struct {
	CmdName          string
	FullCmdName      string
	CmdDepth         int
	FuncName         string
	CobraCmdName     string
	ShortDescription string
	LongDescription  string
	Options          []*tplOption
	Args             *tplArgs
	Commands         []*tplCommand
	ParentCmdName    string
	HasCtx           bool
	HasErr           bool
}

type tplOption struct {
	Name         string
	LongName     string
	ShortName    string
	Description  string
	DefaultValue string
	FlagType     string
}

type tplArgs struct {
	NoArgs bool
}
