package generator

import (
	"github.com/pinzolo/casee"
)

var supportedFlagTypes = map[string]struct{}{
	"BoolP":          struct{}{},
	"BoolSliceP":     struct{}{},
	"StringP":        struct{}{},
	"StringSliceP":   struct{}{},
	"IntP":           struct{}{},
	"IntSliceP":      struct{}{},
	"Int8P":          struct{}{},
	"Int16P":         struct{}{},
	"Int32P":         struct{}{},
	"Int64P":         struct{}{},
	"Float32P":       struct{}{},
	"Float64P":       struct{}{},
	"UintP":          struct{}{},
	"UintSliceP":     struct{}{},
	"Uint8P":         struct{}{},
	"Uint16P":        struct{}{},
	"Uint32P":        struct{}{},
	"Uint64P":        struct{}{},
	"DurationP":      struct{}{},
	"DurationSliceP": struct{}{},
	"IPP":            struct{}{},
	"IPSliceP":       struct{}{},
	"IPMaskP":        struct{}{},
	"BytesBase64P":   struct{}{},
}

func mapTypeToFlagType(typeName string, isArray bool) (flagType string, err error) {

	if typeName == "IP" || typeName == "IPMask" {
		flagType = typeName
	} else {
		flagType = casee.ToPascalCase(typeName)
	}

	if isArray {
		flagType = flagType + "SliceP"
	} else {
		flagType = flagType + "P"
	}

	if flagType == "ByteSliceP" {
		flagType = "BytesBase64P"
	}

	if _, exists := supportedFlagTypes[flagType]; !exists {
		return "", &ErrUnsupportedFlagType{
			OriginalTypeName: typeName,
			IsArray:          isArray,
			MappedTypeName:   flagType,
		}
	}

	return flagType, err
}
