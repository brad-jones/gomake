package generator

import (
	"go/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("buildCmdTree", func() {

	When("given an empty ast package", func() {
		It("should return an empty command tree", func() {
			Expect(buildCmdTree(&ast.Package{})).To(Equal(cmdTree{}))
		})
	})

	When("given a single function declaration", func() {
		It("should return a single function declaration", func() {
			Expect(buildCmdTree(&ast.Package{
				Name: "main",
				Files: map[string]*ast.File{
					"/.gomake/makefile.go": &ast.File{
						Name: &ast.Ident{
							Name: "main",
						},
						Decls: []ast.Decl{
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "Foo",
								},
								Type: &ast.FuncType{},
							},
						},
					},
				},
			})).To(Equal(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{},
					},
					branch: cmdTree{},
				},
			}))
		})
	})

	When("given multiple function declarations", func() {
		It("should return multiple function declarations", func() {
			Expect(buildCmdTree(&ast.Package{
				Name: "main",
				Files: map[string]*ast.File{
					"/.gomake/makefile.go": &ast.File{
						Name: &ast.Ident{
							Name: "main",
						},
						Decls: []ast.Decl{
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "Foo",
								},
								Type: &ast.FuncType{},
							},
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "Bar",
								},
								Type: &ast.FuncType{},
							},
						},
					},
				},
			})).To(Equal(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{},
					},
					branch: cmdTree{},
				},
				"bar": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Bar",
						},
						Type: &ast.FuncType{},
					},
					branch: cmdTree{},
				},
			}))
		})
	})

	When("given sub cmd function declarations", func() {
		It("should return a nested cmd tree structure", func() {
			Expect(buildCmdTree(&ast.Package{
				Name: "main",
				Files: map[string]*ast.File{
					"/.gomake/makefile.go": &ast.File{
						Name: &ast.Ident{
							Name: "main",
						},
						Decls: []ast.Decl{
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "Foo",
								},
								Type: &ast.FuncType{},
							},
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "FooBar",
								},
								Type: &ast.FuncType{},
							},
						},
					},
				},
			})).To(Equal(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{},
					},
					branch: cmdTree{
						"bar": cmdTreeBranch{
							leaf: &ast.FuncDecl{
								Name: &ast.Ident{
									Name: "FooBar",
								},
								Type: &ast.FuncType{},
							},
							branch: cmdTree{},
						},
					},
				},
			}))
		})
	})

	When("given a single sub cmd function declaration", func() {
		It("should return a nested cmd tree structure and fill in the parent cmd with a noop", func() {
			Expect(buildCmdTree(&ast.Package{
				Name: "main",
				Files: map[string]*ast.File{
					"/.gomake/makefile.go": &ast.File{
						Name: &ast.Ident{
							Name: "main",
						},
						Decls: []ast.Decl{
							&ast.FuncDecl{
								Name: &ast.Ident{
									Name: "FooBar",
								},
								Type: &ast.FuncType{},
							},
						},
					},
				},
			})).To(Equal(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "gomake_noop",
						},
						Type: &ast.FuncType{},
					},
					branch: cmdTree{
						"bar": cmdTreeBranch{
							leaf: &ast.FuncDecl{
								Name: &ast.Ident{
									Name: "FooBar",
								},
								Type: &ast.FuncType{},
							},
							branch: cmdTree{},
						},
					},
				},
			}))
		})
	})

})
