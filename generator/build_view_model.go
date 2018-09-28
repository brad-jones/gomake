package generator

import (
	"go/ast"
)

type tplViewModel struct {
	Version   string
	CacheHash string
	Commands  []*tplCommand
}

type tplCommand struct {
	CmdName          string
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
	ShortName    string
	Description  string
	DefaultValue string
	FlagType     string
}

type tplArgs struct {
	NoArgs bool
}

func buildViewModel(nodes *ast.Package) (*tplViewModel, error) {

	viewModel := &tplViewModel{}
	cmdTree := buildCmdTree(nodes)

	if cmds, err := walkCmdTree(cmdTree, "root"); err == nil {
		viewModel.Commands = cmds
	} else {
		return nil, err
	}

	/*
		TODO: Inject the current git commit hash of the project that contains
		the .gomake folder. As a DevOps guy supporting a bunch of other Developers,
		I find it very useful to be to know exactly which version of the task runner
		someone is using. As there can be many different versions across
		different branches, etc.
	*/
	viewModel.Version = "0.0.0"

	return viewModel, nil
}
