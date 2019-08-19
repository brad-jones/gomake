package generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGenerator(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "TestGenerator")
}

/*
var _ = Describe("Generate", func() {

	When("given an example gomake dir", func() {
		It("should generate a runable cobra cli app", func() {
			cwd, _ := os.Getwd()
			dir := filepath.Join(cwd, "..", "example", ".gomake")
			err := generator.Generate(dir)
			Expect(err).ToNot(HaveOccurred())

			out, err := exec.Command("go", "run", dir).Output()
			Expect(err).ToNot(HaveOccurred())
			Expect(string(out)).To(ContainSubstring("Makefile written in golang"))
		})
	})

})

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
					HasErr:           false,
				},
				&tplCommand{
					CmdName:          "bar",
					FullCmdName:      "bar",
					CmdDepth:         0,
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
					FullCmdName:      "foo",
					CmdDepth:         0,
					FuncName:         "Foo",
					CobraCmdName:     "foo",
					ShortDescription: "",
					LongDescription:  "",
					Options:          []*tplOption{},
					Args:             &tplArgs{NoArgs: true},
					Commands: []*tplCommand{
						&tplCommand{
							CmdName:          "bar",
							FullCmdName:      "foo bar",
							CmdDepth:         1,
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

var _ = Describe("getDefaultValueForType", func() {

	When("given an unrecognised type", func() {
		It("should return nil", func() {
			Expect(getDefaultValueForType("foo", false)).To(Equal("nil"))
		})
	})

	When("given an array type", func() {
		It("should return nil", func() {
			Expect(getDefaultValueForType("foo", true)).To(Equal("nil"))
		})
	})

	When("given a string type", func() {
		It("should return an empty string", func() {
			Expect(getDefaultValueForType("string", false)).To(Equal(`""`))
		})
	})

	When("given a bool type", func() {
		It("should return false", func() {
			Expect(getDefaultValueForType("bool", false)).To(Equal("false"))
		})
	})

	When("given a number type", func() {
		It("should return zero", func() {
			Expect(getDefaultValueForType("int", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int16", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("int64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint8", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint16", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("uint64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("float32", false)).To(Equal("0"))
			Expect(getDefaultValueForType("float64", false)).To(Equal("0"))
			Expect(getDefaultValueForType("Duration", false)).To(Equal("0"))
		})
	})

})

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


var _ = Describe("mapTypeToFlagType", func() {

	When("given bool", func() {
		It("should return BoolP", func() {
			result, err := mapTypeToFlagType("bool", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BoolP"))
		})
	})

	When("given bool and is array", func() {
		It("should return BoolSliceP", func() {
			result, err := mapTypeToFlagType("bool", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BoolSliceP"))
		})
	})

	When("given string", func() {
		It("should return StringP", func() {
			result, err := mapTypeToFlagType("string", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("StringP"))
		})
	})

	When("given string and is array", func() {
		It("should return StringSliceP", func() {
			result, err := mapTypeToFlagType("string", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("StringSliceP"))
		})
	})

	When("given int", func() {
		It("should return IntP", func() {
			result, err := mapTypeToFlagType("int", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IntP"))
		})
	})

	When("given int and is array", func() {
		It("should return IntSliceP", func() {
			result, err := mapTypeToFlagType("int", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IntSliceP"))
		})
	})

	When("given int8", func() {
		It("should return Int8P", func() {
			result, err := mapTypeToFlagType("int8", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int8P"))
		})
	})

	When("given int16", func() {
		It("should return Int16P", func() {
			result, err := mapTypeToFlagType("int16", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int16P"))
		})
	})

	When("given int32", func() {
		It("should return Int32P", func() {
			result, err := mapTypeToFlagType("int32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int32P"))
		})
	})

	When("given int64", func() {
		It("should return Int64P", func() {
			result, err := mapTypeToFlagType("int64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Int64P"))
		})
	})

	When("given float32", func() {
		It("should return Float32P", func() {
			result, err := mapTypeToFlagType("float32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Float32P"))
		})
	})

	When("given float64", func() {
		It("should return Float64P", func() {
			result, err := mapTypeToFlagType("float64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Float64P"))
		})
	})

	When("given uint", func() {
		It("should return UintP", func() {
			result, err := mapTypeToFlagType("uint", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("UintP"))
		})
	})

	When("given uint and is array", func() {
		It("should return UintSliceP", func() {
			result, err := mapTypeToFlagType("uint", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("UintSliceP"))
		})
	})

	When("given uint8", func() {
		It("should return Uint8P", func() {
			result, err := mapTypeToFlagType("uint8", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint8P"))
		})
	})

	When("given uint16", func() {
		It("should return Uint16P", func() {
			result, err := mapTypeToFlagType("uint16", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint16P"))
		})
	})

	When("given uint32", func() {
		It("should return Uint32P", func() {
			result, err := mapTypeToFlagType("uint32", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint32P"))
		})
	})

	When("given uint64", func() {
		It("should return Uint64P", func() {
			result, err := mapTypeToFlagType("uint64", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("Uint64P"))
		})
	})

	When("given Duration", func() {
		It("should return DurationP", func() {
			result, err := mapTypeToFlagType("Duration", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("DurationP"))
		})
	})

	When("given Duration and is array", func() {
		It("should return DurationSliceP", func() {
			result, err := mapTypeToFlagType("Duration", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("DurationSliceP"))
		})
	})

	When("given IP", func() {
		It("should return IPP", func() {
			result, err := mapTypeToFlagType("IP", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPP"))
		})
	})

	When("given IP and is array", func() {
		It("should return IPSliceP", func() {
			result, err := mapTypeToFlagType("IP", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPSliceP"))
		})
	})

	When("given IPMask", func() {
		It("should return IPMaskP", func() {
			result, err := mapTypeToFlagType("IPMask", false)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("IPMaskP"))
		})
	})

	When("given byte and is array", func() {
		It("should return BytesBase64P", func() {
			result, err := mapTypeToFlagType("byte", true)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("BytesBase64P"))
		})
	})

	When("given foo an unsupported type", func() {
		It("should return ErrUnsupportedFlagType", func() {
			_, err := mapTypeToFlagType("foo", false)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("gomake: the flag type 'FooP' is unsupported by gomake and/or cobra"))
		})
	})

})

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

*/
