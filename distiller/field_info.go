package distiller

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"regexp"
	"strings"
)

type FieldLayout int

const (
	LayoutSingle FieldLayout = iota // The field is a single element.
	LayoutArray                     // The field is an array or slice of elements.
	LayoutMap                       // The field is a map of elements.
)

// FieldInfo holds information about structure field.
type FieldInfo struct {
	Type       types.Type        // Field type, used to compute fully qualified type string.
	Name       string            // Field name.
	Layout     FieldLayout       // Field layout.
	EltType    types.Type        // Field element type, when the field is a slice or map.
	IsEmbedded bool              // True if field is an embedded struct and Name is an empty string.
	Tags       map[string]string // Tags applied to that field as map of name-value key-pairs.
	Doc        string            // Documentation content if present.
}

// tagRegexp defines a regex to extract tags names and values.
var tagRegexp = regexp.MustCompile(`(\w+):"((?:[^"\\]|\\.)*)"`)

// NewFieldInfo creates new field information object from given abstract syntax tree field and package.
// Terminates the process with a fatal error if multiple names are specified for the same field.
func NewFieldInfo(field *ast.Field, pkg *packages.Package) *FieldInfo {
	f := &FieldInfo{Layout: LayoutSingle, EltType: nil}
	if len(field.Names) == 1 {
		f.Name = field.Names[0].Name
	} else if field.Names == nil {
		// Embedded field.
		f.IsEmbedded = true
	} else {
		log.Fatalf("Unsupported multiple names.")
	}

	f.Type = pkg.TypesInfo.Types[field.Type].Type
	switch fieldType := field.Type.(type) {
	case *ast.ArrayType:
		// In case of array get the type of single element.
		f.EltType = pkg.TypesInfo.Types[fieldType.Elt].Type
		f.Layout = LayoutArray

	case *ast.MapType:
		// In case of map get the type of value.
		f.EltType = pkg.TypesInfo.Types[fieldType.Value].Type
		f.Layout = LayoutMap
	}

	// Parse defined tags populating FieldInfo.Tags map.
	if field.Tag != nil {
		f.Tags = make(map[string]string)
		tags := tagRegexp.FindAllStringSubmatch(strings.Trim(field.Tag.Value, "` "), -1)
		for _, tag := range tags {
			tagValue := ""
			if len(tag) == 3 {
				tagValue = tag[2]
			}
			f.Tags[tag[1]] = tagValue
		}
	}

	// Merge documentation and comment.
	f.Doc = field.Doc.Text() + field.Comment.Text()
	return f
}

func (f *FieldInfo) String() string {
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nLayout: %v\nElement type: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type.String(), f.Name, f.Layout, f.EltType,
		f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}

// FormatDoc formats the field documentation indenting it with passed indent string.
func (f *FieldInfo) FormatDoc(indent string, renderType bool) string {
	doc := f.Doc

	// Check if the type is used to define typed constants.
	consts := LookupTypedConsts(f.Type.String())
	if consts != nil {
		// Display allowed values for defined constants below the field documentation.
		doc += "Allowed values:\n"

		constLen := 0
		valueLen := 0
		for _, info := range consts {
			if len(info.Name) > constLen {
				constLen = len(info.Name)
			}
			if len(info.Value) > valueLen {
				valueLen = len(info.Value)
			}
		}

		for _, info := range consts {
			doc += fmt.Sprintf("%-*s = %*v  %s\n", constLen, info.Name, valueLen, info.Value, info.InlineDoc())
		}
	}

	// Indent the documentation.
	commentPrefix := indent + "// "
	d := strings.ReplaceAll(doc, "\n", "\n"+commentPrefix)
	if len(d) > 0 {
		d = d[:len(d)-len(commentPrefix)]
	} else {
		d = "\n"
	}

	typeName := ""

	if renderType {
		typeName = f.Type.String()
		if lastSlash := strings.LastIndex(typeName, "/"); lastSlash >= 0 {
			typeName = typeName[lastSlash+1:]
			// The square brackets at the beginning of the typeName are trimmed out, so must be re-added.
			if f.Layout == LayoutArray {
				typeName = "[]" + typeName
			}
		}

		if d != "\n" {
			typeName += " - "
		}

		d = typeName + d
	}

	if d == "\n" {
		return ""
	}

	return commentPrefix + d
}
