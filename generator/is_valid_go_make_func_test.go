package generator

import (
	"go/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("isValidGoMakeFunc", func() {

	When("given an un-exported function", func() {
		It("should return false", func() {
			Expect(isValidGoMakeFunc(&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "foo",
				},
				Type: &ast.FuncType{},
			})).To(Equal(false))
		})
	})

	When("given an exported function", func() {
		It("should return true", func() {
			Expect(isValidGoMakeFunc(&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "Foo",
				},
				Type: &ast.FuncType{},
			})).To(Equal(true))
		})
	})

	When("given a function with many return values", func() {
		It("should return false", func() {
			Expect(isValidGoMakeFunc(&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "Foo",
				},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{},
							&ast.Field{},
						},
					},
				},
			})).To(Equal(false))
		})
	})

	When("given a function with a single return value that is not an error", func() {
		It("should return false", func() {
			Expect(isValidGoMakeFunc(&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "Foo",
				},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{Type: &ast.Ident{Name: "foo"}},
						},
					},
				},
			})).To(Equal(false))
		})
	})

	When("given a function with a single return value that is an error", func() {
		It("should return true", func() {
			Expect(isValidGoMakeFunc(&ast.FuncDecl{
				Name: &ast.Ident{
					Name: "Foo",
				},
				Type: &ast.FuncType{
					Results: &ast.FieldList{
						List: []*ast.Field{
							&ast.Field{Type: &ast.Ident{Name: "error"}},
						},
					},
				},
			})).To(Equal(true))
		})
	})

})
