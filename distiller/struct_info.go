package distiller

import (
	"fmt"
	"github.com/marco-sacchi/go2jsonc/ordered"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"strings"
)

// StructInfo holds information about a struct.
type StructInfo struct {
	Package *packages.Package // Package of which the structure is part
	Name    string            // Struct name.
	Doc     string            // Documentation content if present.
	Fields  []*FieldInfo      // Struct fields.

	// If a function with following signature:
	//
	//     func <StructName>Defaults() *StructName
	//
	// is found on the same package where the struct is declared
	// it will be parsed to extract default values for all fields.
	Defaults map[string]interface{} // Map of defaults values for struct fields.
}

func (s *StructInfo) String() string {
	return fmt.Sprintf(
		"Package: %s\nName: \"%s\"\nFieldCount: %v\nFields:\n%s\nDoc: \"%v\"\nDefaults: %+v",
		s.Package.PkgPath, s.Name, len(s.Fields), s.Fields,
		strings.ReplaceAll(s.Doc, "\n", "\\n"), s.Defaults,
	)
}

// NewStructInfo creates a new struct information object from given abstract syntax tree type spec
// and types info read on package loading.
func NewStructInfo(genDecl *ast.GenDecl, pkg *packages.Package) *StructInfo {
	typeSpec := genDecl.Specs[0].(*ast.TypeSpec)
	structType := typeSpec.Type.(*ast.StructType)

	info := &StructInfo{
		Package: pkg,
		Name:    typeSpec.Name.Name,
		Doc:     genDecl.Doc.Text() + typeSpec.Doc.Text() + typeSpec.Comment.Text(),
	}

	for _, field := range structType.Fields.List {
		f := NewFieldInfo(field, info.Package)
		info.Fields = append(info.Fields, f)
	}

	return info
}

// FormatDoc formats the struct documentation indenting it with passed indent string.
func (s *StructInfo) FormatDoc(indent string) string {
	commentPrefix := indent + "// "
	d := strings.ReplaceAll(s.Doc, "\n", "\n"+commentPrefix)
	return commentPrefix + d[:len(d)-len(commentPrefix)]
}

// ParseDefaultsMethod parses the Defaults method of this struct populating Defaults map.
func (s *StructInfo) ParseDefaultsMethod() error {
	typePath := s.Package.PkgPath + "." + s.Name
	for _, astFile := range s.Package.Syntax {
		for _, decl := range astFile.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok || funcDecl.Name.Name != s.Name+"Defaults" {
				continue
			}

			if funcDecl.Recv != nil || funcDecl.Type.Results.NumFields() != 1 ||
				s.Package.TypesInfo.Types[funcDecl.Type.Results.List[0].Type].Type.String() != "*"+typePath {
				return fmt.Errorf("invalid defaults method signature.\n"+
					"expected: func %sDefaults() *%s\n"+
					"got:      func %s() %s",
					s.Name, typePath, funcDecl.Name.Name,
					s.Package.TypesInfo.Types[funcDecl.Type.Results.List[0].Type].Type.String())
			}

			var defaults interface{}
			ast.Inspect(funcDecl, func(node ast.Node) bool {
				switch n := node.(type) {
				case *ast.CompositeLit:
					var ident *ast.Ident
					if ident, ok = n.Type.(*ast.Ident); ok && ident.Name == s.Name {
						defaults = s.parseDefaultsMethodBody(n)
						// Stop traversing.
						return false
					}
				}

				// Continue traversing.
				return true
			})

			s.Defaults = defaults.(map[string]interface{})
		}
	}

	return nil
}

// ParseDefaultsMethod parses recursively the composite literals of Defaults method. It returns a map of
// fields names-values or an array of values.
func (s *StructInfo) parseDefaultsMethodBody(lit *ast.CompositeLit) interface{} {
	var values interface{}
	isOrdered := false
	switch lit.Type.(type) {
	case *ast.ArrayType:
		values = make([]interface{}, 0)

	case *ast.MapType:
		// Uses an ordered map to ensure consistent and sorted output according
		// to the order in which the default values are defined.
		isOrdered = true
		values = ordered.NewMap()

	default:
		values = make(map[string]interface{})
	}

	for _, elt := range lit.Elts {
		switch el := elt.(type) {
		case *ast.KeyValueExpr:
			// values is a key-value pair.
			var key string
			switch k := el.Key.(type) {
			case *ast.Ident:
				key = fmt.Sprintf("%s", k.Name)

			case *ast.BasicLit:
				key = fmt.Sprintf("%s", k.Value)
			}

			var value interface{}
			switch val := el.Value.(type) {
			case *ast.CompositeLit:
				value = s.parseDefaultsMethodBody(val)

			default:
				value = s.Package.TypesInfo.Types[val].Value
			}

			if isOrdered {
				values.(*ordered.Map).Append(key, value)
			} else {
				values.(map[string]interface{})[key] = value
			}

		case *ast.CompositeLit:
			// values is an array of interfaces.
			values = append(values.([]interface{}), s.parseDefaultsMethodBody(el))

		default:
			// values is an array of interfaces.
			values = append(values.([]interface{}), s.Package.TypesInfo.Types[el].Value)
		}
	}

	return values
}
