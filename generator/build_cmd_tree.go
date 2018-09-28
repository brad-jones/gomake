package generator

import (
	"go/ast"
	"strings"

	"github.com/pinzolo/casee"
)

type cmdTree map[string]cmdTreeBranch

type cmdTreeBranch struct {
	leaf   *ast.FuncDecl
	branch cmdTree
}

func buildCmdTree(nodes *ast.Package) cmdTree {

	tree := cmdTree{}

	ast.Inspect(nodes, func(node ast.Node) bool {
		switch funcDecl := node.(type) {
		case *ast.FuncDecl:
			if isValidGoMakeFunc(funcDecl) {
				root := tree
				preserveUnderscores := strings.Replace(funcDecl.Name.Name, "_", ":", -1)
				snakeCase := casee.ToSnakeCase(preserveUnderscores)
				convertPreservedUnderscoresToHyphens := strings.Replace(snakeCase, "_:_", "-", -1)
				for _, cmd := range strings.Split(convertPreservedUnderscoresToHyphens, "_") {
					if _, ok := root[cmd]; !ok {
						root[cmd] = cmdTreeBranch{
							leaf:   funcDecl,
							branch: cmdTree{},
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
