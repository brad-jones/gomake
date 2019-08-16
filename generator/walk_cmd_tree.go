package generator

import (
	"regexp"
	"strings"

	"github.com/pinzolo/casee"
)

var regexpOptDoc = regexp.MustCompile("--(.*?),?:")

func walkCmdTree(tree cmdTree, packageName, parentCmdName, fullCmdName string, depth int) ([]*tplCommand, error) {

	cmds := []*tplCommand{}

	for cmdName, branch := range tree {

		// Add a new command
		cmd := &tplCommand{}
		cmd.CmdName = cmdName
		cmd.FullCmdName = strings.TrimSpace(fullCmdName + " " + cmdName)
		cmd.CmdDepth = depth

		if branch.leaf.Name.Name != "gomake_noop" {
			cmd.FuncName = packageName + "." + branch.leaf.Name.Name
		} else {
			cmd.FuncName = branch.leaf.Name.Name
		}

		cmd.CobraCmdName = casee.ToCamelCase(cmdName)
		cmd.ParentCmdName = casee.ToCamelCase(parentCmdName + "Cmd")
		cmds = append(cmds, cmd)

		// Add any sub commands
		if subCmds, err := walkCmdTree(branch.branch, packageName, cmdName, cmd.FullCmdName, depth+1); err == nil {
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
		optDefaultValues := map[string]string{}
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
						// --baz, $BAZ: yes defaults can be provided via the env
						// --abc, -a, $ABC: also supported
						// --qux, "a static default": this should also work
						// They may only span a single line
						parts := strings.Split(regexpOptDoc.FindStringSubmatch(line)[1], ",")
						value := strings.TrimSpace(regexpOptDoc.ReplaceAllString(line, ""))
						optDocs[parts[0]] = value
						if len(parts) > 1 {
							v := strings.TrimSpace(parts[1])
							if strings.HasPrefix(v, "-") {
								optShortNames[parts[0]] = strings.TrimPrefix(v, "-")
							} else {
								optDefaultValues[parts[0]] = v
							}
							if len(parts) > 2 {
								optDefaultValues[parts[0]] = strings.TrimSpace(parts[2])
							}
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
					description := optDocs[optNameChain]
					optDefaultValue := defaultValue
					if val, ok := optDefaultValues[optNameChain]; ok {
						if strings.HasPrefix(val, "$") {
							description = "env: " + val + " - " + description
							val = "os.Getenv(\"" + strings.TrimPrefix(val, "$") + "\")"
						}
						optDefaultValue = val
					}
					optDefaultValue = strings.Replace(optDefaultValue, "\\\"", "\"", -1)
					cmd.Options = append(cmd.Options, &tplOption{
						Name:         name.Name,
						LongName:     optNameChain,
						ShortName:    optShortNames[optNameChain],
						Description:  description,
						DefaultValue: optDefaultValue,
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
