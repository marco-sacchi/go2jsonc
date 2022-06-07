package main

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"regexp"
	"strings"
)

// FieldInfo holds information about structure field.
type FieldInfo struct {
	Type       types.Type        // Field type, used to compute fully qualified type string.
	Name       string            // Field name.
	IsArray    bool              // True if field type is an array.
	IsEmbedded bool              // True if field is an embedded struct and Name is an empty string.
	Tags       map[string]string // Tags applied to that field as map of name-value key-pairs.
	Doc        string            // Documentation content if present.
}

// tagRegexp defines a regex to extract tags names and values.
var tagRegexp = regexp.MustCompile(`(\w+):"((?:[^"\\]|\\.)*)"`)

// NewFieldInfo creates new field information object from given abstract syntax tree field and package.
// Terminates the process with a fatal error if multiple names are specified for the same field.
func NewFieldInfo(field *ast.Field, pkg *packages.Package) *FieldInfo {
	f := &FieldInfo{IsArray: false}
	if len(field.Names) == 1 {
		f.Name = field.Names[0].Name
	} else if field.Names == nil {
		// Embedded field.
		f.IsEmbedded = true
	} else {
		log.Fatalf("Unsupported multiple names.")
	}

	f.Type = pkg.TypesInfo.Types[field.Type].Type
	switch field.Type.(type) {
	case *ast.ArrayType:
		// In case of array get the type of a single element.
		f.Type = pkg.TypesInfo.Types[field.Type.(*ast.ArrayType).Elt].Type
		f.IsArray = true
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
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nIsArray: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type.String(), f.Name, f.IsArray, f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}

// FormatDoc formats the field documentation indenting it with passed indent string.
func (f *FieldInfo) FormatDoc(indent string) string {
	doc := f.Doc

	// Check if the type is used to define typed constants.
	consts := lookupTypedConsts(f.Type.String())
	if consts != nil {
		// Display allowed values for defined constants below the field documentation.
		doc += "Allowed values:\n"
		for _, info := range consts {
			doc += fmt.Sprintf("%s = %v\n", info.Name, info.Value)
		}
	}

	// Indent the documentation.
	commentPrefix := indent + "// "
	d := strings.ReplaceAll(doc, "\n", "\n"+commentPrefix)
	if len(d) > 0 {
		d = " - " + d[:len(d)-len(commentPrefix)]
	} else {
		d = "\n"
	}

	typeName := f.Type.String()
	if f.IsArray {
		typeName += "[]"
	}

	if lastSlash := strings.LastIndex(typeName, "/"); lastSlash >= 0 {
		typeName = typeName[lastSlash+1:]
	}

	return commentPrefix + typeName + d
}
