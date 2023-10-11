// Package go2jsonc generates the jsonc code from the information extracted
// from the AST via the distiller package.
package go2jsonc

import (
	"fmt"
	"github.com/marco-sacchi/go2jsonc/distiller"
	"github.com/marco-sacchi/go2jsonc/ordered"
	"go/constant"
	"go/types"
	"log"
	"strings"
)

// DocTypesMode defines rendering modes for field types in JSONC comments.
type DocTypesMode int

const (
	NotStructFields DocTypesMode = 1 << iota // Don't show type on struct fields.
	NotArrayFields                           // Don't show type on array or slice fields.
	NotMapFields                             // Don't show type on map fields.
	AllFields       DocTypesMode = 0         // Show types on all fields (default).
)

var docTypesMode = AllFields

// Generate generates JSONC indented code for given package dir and type name.
// mode controls the rendering of field types in JSONC comments.
func Generate(dir, typeName string, mode DocTypesMode) (string, error) {
	pkgInfo, err := distiller.NewPackageInfo(dir, typeName)
	if err != nil {
		return "", err
	}

	s := distiller.LookupStruct(pkgInfo.Package.PkgPath + "." + typeName)
	if s == nil {
		return "", fmt.Errorf("cannot find struct %s in package %s", typeName, pkgInfo.Package.Name)
	}

	docTypesMode = mode

	var code string
	code, err = renderStruct(s, s.Defaults, "", false, nil)
	if err != nil {
		log.Fatal(err)
	}

	return code, nil
}

// renderStruct renders JSONC indented code for specified struct and all nested or embedded ones recursively.
func renderStruct(info *distiller.StructInfo, defaults interface{}, indent string,
	embedded bool, parentShadowing []string) (string, error) {
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

	comma := ""
	block_spacing := false
	for i, field := range info.Fields {
		name := field.Name

		// This field will be shadowed by another one, so skip it.
		if (!field.IsEmbedded && lastIndexOf(shadowing, name) > i) ||
			(embedded && lastIndexOf(parentShadowing, name) != -1) {
			continue
		}

		builder.WriteString(comma)

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

		consts := distiller.LookupTypedConsts(field.Type.String())

		renderType := true

		// No default defined for this field.
		if !ok {
			if consts != nil {
				value = consts[0].Value
			} else {
				value = typeZero(field)
			}
		} else {
			var err error
			switch field.Layout {
			case distiller.LayoutSingle:
				if _, ok = field.Type.(*types.Named); ok && consts == nil {
					subInfo := distiller.LookupStruct(field.Type.String())
					if subInfo == nil {
						return "", fmt.Errorf("cannot lookup structure %s", field.Type.String())
					}

					renderType = (docTypesMode & NotStructFields) == 0

					value, err = renderStruct(subInfo, value, indent, field.IsEmbedded, shadowing[i:])
					if err != nil {
						return "", err
					}
				}

				// No special handling required for basic types.

			case distiller.LayoutArray:
				renderType = (docTypesMode & NotArrayFields) == 0
				value, err = renderArray(field, value.([]interface{}), indent)

			case distiller.LayoutMap:
				renderType = (docTypesMode & NotMapFields) == 0
				value, err = renderMap(field, value.(*ordered.Map), indent)
			}

			if err != nil {
				return "", err
			}
		}

		if field.IsEmbedded {
			builder.WriteString(fmt.Sprintf("%v", value))
		} else {
			doc := field.FormatDoc(indent, renderType)
			if doc != "" {
				// Adds a blank line when the comment block is present.
				if !block_spacing && (comma != "") {
					builder.WriteString("\n")
				}
				block_spacing = true
			} else {
				block_spacing = false
			}

			builder.WriteString(doc)
			builder.WriteString(fmt.Sprintf("%s\"%s\": %v", indent, name, value))
		}

		comma = ",\n"
		if block_spacing {
			comma += "\n"
		}
	}

	if !embedded {
		if comma != "" {
			builder.WriteString("\n")
		}

		builder.WriteString(indent[:len(indent)-1] + "}")
	}

	return builder.String(), nil
}

// renderArray renders slice or array fields.
func renderArray(field *distiller.FieldInfo, value []interface{}, indent string) (string, error) {
	if len(value) == 0 {
		return "[]", nil
	}

	eltsIdent := indent + "\t"
	code := "[\n"
	for _, elt := range value {
		literal, err := renderElement(field.EltType, elt, eltsIdent)
		if err != nil {
			return "", err
		}

		code += eltsIdent + literal + ",\n"
	}
	code = strings.TrimRight(code, ",\n") + "\n" + indent + "]"

	return code, nil
}

// renderMap renders map fields.
func renderMap(field *distiller.FieldInfo, value *ordered.Map, indent string) (string, error) {
	if field.IsEmbedded == true {
		return "", fmt.Errorf("field of slice or map type cannot be embedded")
	}

	if value.Len() == 0 {
		return "{}", nil
	}

	eltsIndent := indent + "\t"
	code := "{\n"

	var err error
	value.Iterate(func(key string, elt interface{}) bool {
		var literal string
		literal, err = renderElement(field.EltType, elt, eltsIndent)
		if err != nil {
			return false
		}

		code += eltsIndent + fmt.Sprintf("%s: %s", key, literal) + ",\n"
		return true
	})

	if err != nil {
		return "", err
	}

	code = strings.TrimRight(code, ",\n") + "\n" + indent + "}"

	return code, nil
}

// renderElement renders an element value of a slice, array or map.
func renderElement(itemType types.Type, item interface{}, indent string) (string, error) {
	_, ok := itemType.(*types.Basic)
	if ok || distiller.LookupTypedConsts(itemType.String()) != nil {
		return fmt.Sprintf("%v", item), nil
	}

	subInfo := distiller.LookupStruct(itemType.String())
	if subInfo == nil {
		return "", fmt.Errorf("cannot lookup structure %s", itemType.String())
	}

	return renderStruct(subInfo, item, indent, false, nil)
}

// lastIndexOf returns the last slice index of specified value.
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

// typeZero return the default uninitialized value for specified field.
func typeZero(field *distiller.FieldInfo) interface{} {
	var value interface{}
	if field.Layout == distiller.LayoutArray {
		value = make([]interface{}, 0)
		return value
	} else if field.Layout == distiller.LayoutMap {
		value = make(map[interface{}]interface{}, 0)
		return value
	}

	fieldType := types.Default(field.Type)
	switch t := fieldType.(type) {
	case *types.Basic:
		switch t.Kind() {
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
			log.Fatalf("Unhandled default value for type %v", fieldType.String())
		}
	}

	return value
}
