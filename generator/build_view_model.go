package generator

import (
	"go/ast"
	"strings"

	git "gopkg.in/src-d/go-git.v4"
)

type tplViewModel struct {
	Use       string
	Short     string
	Version   string
	CacheHash string
	Commands  []*tplCommand
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

func buildViewModel(nodes *ast.Package) (*tplViewModel, error) {

	viewModel := &tplViewModel{}
	cmdTree := buildCmdTree(nodes)

	if cmds, err := walkCmdTree(cmdTree, "root", "", 0); err == nil {
		viewModel.Commands = cmds
	} else {
		return nil, err
	}

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

	return viewModel, nil
}
