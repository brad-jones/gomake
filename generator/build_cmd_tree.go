package generator

import (
	"go/ast"
	"regexp"
	"strings"

	"github.com/pinzolo/casee"
)

type cmdTree map[string]cmdTreeBranch

type cmdTreeBranch struct {
	leaf   *ast.FuncDecl
	branch cmdTree
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

func buildCmdTree(nodes *ast.Package) cmdTree {

	tree := cmdTree{}

	ast.Inspect(nodes, func(node ast.Node) bool {
		switch funcDecl := node.(type) {
		case *ast.FuncDecl:
			if isValidGoMakeFunc(funcDecl) {
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
