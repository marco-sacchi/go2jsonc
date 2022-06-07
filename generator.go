package main

import (
	"fmt"
	"go/constant"
	"go/types"
	"log"
	"strings"
)

// GenerateJSONC generates JSONC indented code.
func GenerateJSONC(dir, typeName string) (string, error) {
	pkgInfo, err := NewPackageInfo(dir, typeName)
	if err != nil {
		return "", err
	}

	s, ok := pkgInfo.Structs[pkgInfo.Package.PkgPath+"."+typeName]
	if !ok {
		return "", fmt.Errorf("cannot find struct %s in package %s", typeName, pkgInfo.Package.Name)
	}

	var code string
	code, err = renderStruct(s, s.Defaults, "", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return code, nil
}

// renderStruct renders JSONC indented code for specified struct and all nested or embedded ones recursively.
func renderStruct(info *StructInfo, defaults interface{}, indent string, embedded bool, parentShadowing []string) (string, error) {
	var builder strings.Builder

	if !embedded {
		builder.WriteString("{\n")
		indent += "\t"
	}

	var shadowing []string
	for _, field := range info.Fields {
		if !field.IsEmbedded {
			shadowing = append(shadowing, field.Name)
		}
	}

	for i, field := range info.Fields {
		name := field.Name

		// This field will be shadowed by another one, so skip it.
		if (!field.IsEmbedded && lastIndexOf(shadowing, name) > i) ||
			(embedded && lastIndexOf(parentShadowing, name) != -1) {
			continue
		}

		if jsonName, ok := field.Tags["json"]; ok == true {
			name = jsonName
		}

		key := field.Name
		if field.IsEmbedded {
			key = field.Type.String()
			if pathEnd := strings.LastIndex(key, "/"); pathEnd >= 0 {
				key = key[pathEnd+strings.Index(key[pathEnd+1:], ".")+2:]
			}
		}

		var value interface{}
		ok := false
		if defaults != nil {
			value, ok = defaults.(map[string]interface{})[key]
		}

		consts := lookupTypedConsts(field.Type.String())

		if !ok {
			if consts != nil {
				value = consts[0].Value
			} else {
				value = typeZero(field)
			}
		}

		if _, ok = field.Type.(*types.Basic); ok || consts != nil {
			if field.IsArray {
				arrayIndent := indent + "\t"
				array := "[\n"
				for _, item := range value.([]interface{}) {
					array += arrayIndent + fmt.Sprintf("%v", item) + ",\n"
				}
				value = strings.TrimRight(array, ",\n") + "\n" + indent + "]"
				if value == "[\n"+indent+"]" {
					value = "[]"
				}
			}
		} else {
			subInfo := lookupStruct(field.Type.String())
			if subInfo == nil {
				return "", fmt.Errorf("cannot lookup structure %s", field.Type.String())
			}

			var err error
			if field.IsArray {
				arrayIndent := indent + "\t"
				array := "[\n"
				var itemString string
				for _, item := range value.([]interface{}) {
					itemString, err = renderStruct(subInfo, item, arrayIndent, field.IsEmbedded, shadowing[i:])
					if err != nil {
						return "", err
					}
					array += arrayIndent + itemString + ",\n"
				}
				value = strings.TrimRight(array, ",\n") + "\n" + indent + "]"
				if value == "[\n"+indent+"]" {
					value = "[]"
				}
			} else {
				value, err = renderStruct(subInfo, value, indent, field.IsEmbedded, shadowing[i:])
				if err != nil {
					return "", err
				}
			}
		}

		if field.IsEmbedded {
			builder.WriteString(fmt.Sprintf("%v", value))
		} else {
			builder.WriteString(field.FormatDoc(indent))
			builder.WriteString(fmt.Sprintf("%s\"%s\": %v,\n", indent, name, value))
		}
	}

	if !embedded {
		indent = indent[:len(indent)-1]
		builder.WriteString(indent + "}")
	}

	return builder.String(), nil
}

func lastIndexOf(slice []string, value string) int {
	if slice != nil {
		for i := len(slice) - 1; i >= 0; i-- {
			if slice[i] == value {
				return i
			}
		}
	}

	return -1
}

func typeZero(field *FieldInfo) interface{} {
	var value interface{}
	if field.IsArray {
		value = make([]interface{}, 0)
		return value
	}

	t := types.Default(field.Type)
	switch t.(type) {
	case *types.Basic:
		switch t.(*types.Basic).Kind() {
		case types.Bool:
			value = *new(bool)
		case types.Int:
			value = *new(int)
		case types.Int8:
			value = *new(int8)
		case types.Int16:
			value = *new(int16)
		case types.Int32:
			value = *new(int32)
		case types.Int64:
			value = *new(int64)
		case types.Uint:
			value = *new(uint)
		case types.Uint8:
			value = *new(uint8)
		case types.Uint16:
			value = *new(uint16)
		case types.Uint32:
			value = *new(uint32)
		case types.Uint64:
			value = *new(uint64)
		case types.Uintptr:
			value = *new(uintptr)
		case types.Float32:
			value = *new(float32)
		case types.Float64:
			value = *new(float64)
		case types.Complex64:
			value = *new(complex64)
		case types.Complex128:
			value = *new(complex128)
		case types.String:
			value = constant.MakeString("")
		default:
			log.Fatalf("Unhandled default value for type %v", t.String())
		}
	}

	return value
}
