package generator

import (
	"go/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("getTypeFromAstNode", func() {

	When("given an ast.Ident", func() {
		It("should return its name and not be an array or positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.Ident{
					Name: "foo",
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(false))
			Expect(isPositional).To(Equal(false))
		})
	})

	When("given an ast.SelectorExpr", func() {
		It("should return its name and not be an array or positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.SelectorExpr{
					Sel: &ast.Ident{
						Name: "foo",
					},
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(false))
			Expect(isPositional).To(Equal(false))
		})
	})

	When("given an ast.ArrayType[ast.Ident]", func() {
		It("should return its name and be an array but not be positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.ArrayType{
					Elt: &ast.Ident{
						Name: "foo",
					},
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(true))
			Expect(isPositional).To(Equal(false))
		})
	})

	When("given an ast.ArrayType[ast.SelectorExpr]", func() {
		It("should return its name and be an array but not be positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.ArrayType{
					Elt: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "foo",
						},
					},
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(true))
			Expect(isPositional).To(Equal(false))
		})
	})

	When("given an ast.Ellipsis[ast.Ident]", func() {
		It("should return its name and be an array and be positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.Ellipsis{
					Elt: &ast.Ident{
						Name: "foo",
					},
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(true))
			Expect(isPositional).To(Equal(true))
		})
	})

	When("given an ast.Ellipsis[ast.SelectorExpr]", func() {
		It("should return its name and be an array and be positional", func() {
			typeName, isArray, isPositional := getTypeFromAstNode(&ast.Field{
				Type: &ast.Ellipsis{
					Elt: &ast.SelectorExpr{
						Sel: &ast.Ident{
							Name: "foo",
						},
					},
				},
			})
			Expect(typeName).To(Equal("foo"))
			Expect(isArray).To(Equal(true))
			Expect(isPositional).To(Equal(true))
		})
	})

})
