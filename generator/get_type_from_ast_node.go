package generator

import "go/ast"

func getTypeFromAstNode(node *ast.Field) (typeName string, isArray bool, isPositional bool) {

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
