package generator

func getDefaultValueForType(typeName string, isArray bool) string {

	if isArray {
		return "nil"
	}

	switch typeName {
	case "string":
		return `""`
	case "int",
		"int8",
		"int16",
		"int32",
		"int64",
		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",
		"float32",
		"float64",
		"Duration":
		return "0"
	case "bool":
		return "false"
	}

	return "nil"
}
