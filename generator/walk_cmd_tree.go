package generator

import (
	"regexp"
	"strings"

	"github.com/pinzolo/casee"
)

var regexpOptDoc = regexp.MustCompile("--(.*?),?:")

func walkCmdTree(tree cmdTree, parentCmdName string) ([]*tplCommand, error) {

	cmds := []*tplCommand{}

	for cmdName, branch := range tree {

		// Add a new command
		cmd := &tplCommand{}
		cmd.CmdName = cmdName
		cmd.FuncName = branch.leaf.Name.Name
		cmd.CobraCmdName = casee.ToCamelCase(cmdName)
		cmd.ParentCmdName = casee.ToCamelCase(parentCmdName + "Cmd")
		cmds = append(cmds, cmd)

		// Add any sub commands
		if subCmds, err := walkCmdTree(branch.branch, cmdName); err == nil {
			cmd.Commands = subCmds
		} else {
			return nil, err
		}

		// If the function has a doc comment we will use it to provide
		// the long and short descriptions of the cobra command.
		// The first line being the short description, with anything
		// following making up the remainder of the long description.
		optDocs := map[string]string{}
		optShortNames := map[string]string{}
		if branch.leaf.Doc != nil {
			for k, v := range branch.leaf.Doc.List {
				line := strings.TrimLeft(v.Text, "// ")
				line = strings.Replace(line, "\"", "\\\"", -1)
				if k == 0 {
					cmd.ShortDescription = line
				} else {
					// Lines that start with "--" document command
					// options they should not be included in the
					// overall command documentation
					if !strings.HasPrefix(line, "--") {
						cmd.LongDescription = cmd.LongDescription + "\\n" + line
					} else {
						// Parse option comments, they look like this:
						// --foo, -f: the foo option does xyz
						// --bar: the bar option does abc
						// They may only span a single line
						parts := strings.Split(regexpOptDoc.FindStringSubmatch(line)[1], ",")
						value := strings.TrimSpace(regexpOptDoc.ReplaceAllString(line, ""))
						optDocs[parts[0]] = value
						if len(parts) > 1 {
							optShortNames[parts[0]] = strings.TrimPrefix(strings.TrimSpace(parts[1]), "-")
						}
					}
				}
			}
		}

		cmd.LongDescription = "\\n" + strings.TrimLeft(cmd.LongDescription, "\\n")
		if cmd.LongDescription == "\\n" {
			cmd.LongDescription = ""
		}

		// Function parameters are converted to cobra command options.
		// A final variadic parameter is converted to cobra positional arguments
		// NOTE: Only `...string` is supported
		cmd.Options = []*tplOption{}
		cmd.Args = &tplArgs{NoArgs: true}
		if branch.leaf.Type.Params != nil {
			for _, param := range branch.leaf.Type.Params.List {

				typeName, isArray, isPositional := getTypeFromAstNode(param)

				if typeName == "Context" {
					cmd.HasCtx = true
					continue
				}

				if isPositional {
					if typeName != "string" {
						return nil, &ErrVariadicMustBeString{}
					}
					cmd.Args.NoArgs = false
					continue
				}

				flagType, err := mapTypeToFlagType(typeName, isArray)
				if err != nil {
					return nil, err
				}

				defaultValue := getDefaultValueForType(typeName, isArray)

				for _, name := range param.Names {
					optNameChain := casee.ToChainCase(name.Name)
					cmd.Options = append(cmd.Options, &tplOption{
						Name:         name.Name,
						LongName:     optNameChain,
						ShortName:    optShortNames[optNameChain],
						Description:  optDocs[optNameChain],
						DefaultValue: defaultValue,
						FlagType:     flagType,
					})
				}
			}
		}

		if branch.leaf.Type.Results != nil {
			cmd.HasErr = true
		}
	}

	return cmds, nil
}
