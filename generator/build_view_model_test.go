package generator

import (
	"go/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("buildViewModel", func() {

	When("given an empty ast package", func() {
		It("should return a view model with 0 commands", func() {
			result, err := buildViewModel(&ast.Package{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Commands).To(Equal([]*tplCommand{}))
		})
	})

	When("given an unsuported parameter type", func() {
		It("should return an error", func() {
			_, err := buildViewModel(&ast.Package{
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
								Type: &ast.FuncType{
									Params: &ast.FieldList{
										List: []*ast.Field{
											&ast.Field{
												Type: &ast.Ident{
													Name: "foo",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			})
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: the flag type 'FooP' is unsupported by gomake and/or cobra"))
		})
	})

	When("given a single function declaration", func() {
		It("should return a view model with 1 command", func() {
			result, err := buildViewModel(&ast.Package{
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
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Commands).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           false,
					HasErr:           false,
				},
			}))
		})
	})

	When("given multiple function declarations", func() {
		It("should return a view model with multiple commands", func() {
			result, err := buildViewModel(&ast.Package{
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
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Commands).To(ConsistOf(
				&tplCommand{
					CmdName:          "foo",
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           false,
					HasErr:           false,
				},
				&tplCommand{
					CmdName:          "bar",
					FuncName:         "Bar",
					CobraCmdName:     "bar",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           false,
					HasErr:           false,
				},
			))
		})
	})

	When("given sub cmd function declarations", func() {
		It("should return a view model with nested commands", func() {
			result, err := buildViewModel(&ast.Package{
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
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(result.Commands).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands: []*tplCommand{
						&tplCommand{
							CmdName:          "bar",
							FuncName:         "FooBar",
							CobraCmdName:     "bar",
							ShortDescription: "",
							LongDescription:  "",
							Options:          []*tplOption{},
							Args:             &tplArgs{NoArgs: true},
							Commands:         []*tplCommand{},
							ParentCmdName:    "fooCmd",
							HasCtx:           false,
							HasErr:           false,
						},
					},
					ParentCmdName: "rootCmd",
					HasCtx:        false,
					HasErr:        false,
				},
			}))
		})
	})

})
