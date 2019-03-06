package generator

import (
	"go/ast"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("walkCmdTree", func() {

	When("given an empty cmd tree", func() {
		It("should return no commands", func() {
			result, err := walkCmdTree(cmdTree{}, "root", "", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]*tplCommand{}))
		})
	})

	When("given a single command in the cmd tree", func() {
		It("should return a single command", func() {
			result, err := walkCmdTree(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{Name: "bar"},
										},
										Type: &ast.Ident{Name: "string"},
									},
									&ast.Field{
										Names: []*ast.Ident{
											&ast.Ident{Name: "qux"},
										},
										Type: &ast.Ident{Name: "int"},
									},
								},
							},
						},
						Doc: &ast.CommentGroup{
							List: []*ast.Comment{
								&ast.Comment{
									Text: "// This is short description",
								},
								&ast.Comment{
									Text: "// This is the long description",
								},
								&ast.Comment{
									Text: "// --bar: This is the bar flag description",
								},
								&ast.Comment{
									Text: "// --qux, -q: This is the qux flag description with short flag",
								},
							},
						},
					},
					branch: cmdTree{},
				},
			}, "root", "", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FullCmdName:      "foo",
					CmdDepth:         0,
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "This is short description",
					LongDescription:  "\\nThis is the long description",
					Options: []*tplOption{
						&tplOption{
							Name:         "bar",
							LongName:     "bar",
							ShortName:    "",
							Description:  "This is the bar flag description",
							DefaultValue: "\"\"",
							FlagType:     "StringP",
						},
						&tplOption{
							Name:         "qux",
							LongName:     "qux",
							ShortName:    "q",
							Description:  "This is the qux flag description with short flag",
							DefaultValue: "0",
							FlagType:     "IntP",
						},
					},
					Args:          &tplArgs{NoArgs: true},
					Commands:      []*tplCommand{},
					ParentCmdName: "rootCmd",
					HasCtx:        false,
					HasErr:        false,
				},
			}))
		})
	})

	When("given a command that returns an error", func() {
		It("should return with HasErr set to true", func() {
			result, err := walkCmdTree(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
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
					},
					branch: cmdTree{},
				},
			}, "root", "", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FullCmdName:      "foo",
					CmdDepth:         0,
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           false,
					HasErr:           true,
				},
			}))
		})
	})

	When("given a command that has context", func() {
		It("should return with HasCtx set to true", func() {
			result, err := walkCmdTree(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Type: &ast.SelectorExpr{
											Sel: &ast.Ident{
												Name: "Context",
											},
										},
									},
								},
							},
						},
					},
					branch: cmdTree{},
				},
			}, "root", "", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FullCmdName:      "foo",
					CmdDepth:         0,
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           true,
					HasErr:           false,
				},
			}))
		})
	})

	When("given a command that has a variadic argument", func() {
		It("should return with NoArgs set to false", func() {
			result, err := walkCmdTree(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Type: &ast.Ellipsis{
											Elt: &ast.Ident{
												Name: "string",
											},
										},
									},
								},
							},
						},
					},
					branch: cmdTree{},
				},
			}, "root", "", 0)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal([]*tplCommand{
				&tplCommand{
					CmdName:          "foo",
					FullCmdName:      "foo",
					CmdDepth:         0,
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: false},
					Commands:         []*tplCommand{},
					ParentCmdName:    "rootCmd",
					HasCtx:           false,
					HasErr:           false,
				},
			}))
		})
	})

	When("given a command that has a variadic argument that is not of type string", func() {
		It("should return an error", func() {
			_, err := walkCmdTree(cmdTree{
				"foo": cmdTreeBranch{
					leaf: &ast.FuncDecl{
						Name: &ast.Ident{
							Name: "Foo",
						},
						Type: &ast.FuncType{
							Params: &ast.FieldList{
								List: []*ast.Field{
									&ast.Field{
										Type: &ast.Ellipsis{
											Elt: &ast.Ident{
												Name: "int",
											},
										},
									},
								},
							},
						},
					},
					branch: cmdTree{},
				},
			}, "root", "", 0)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: variadic parameters must of type 'string'"))
		})
	})

	When("given a bad sub command", func() {
		It("should return an error", func() {
			_, err := walkCmdTree(cmdTree{
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
									Name: "Bar",
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
							branch: cmdTree{},
						},
					},
				},
			}, "root", "", 0)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: the flag type 'FooP' is unsupported by gomake and/or cobra"))
		})
	})

})
