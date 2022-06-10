package distiller

import (
	"fmt"
	"go/ast"
	"go/types"
	"golang.org/x/tools/go/packages"
	"strings"
)

// ConstInfo holds information about a typed constant.
type ConstInfo struct {
	Name  string // Constant name.
	Value string // String representation of constant value.
	Doc   string // Constant documentation and comment nodes contents.
}

// NewConstInfo creates new const information object from given abstract syntax tree value spec and package.
func NewConstInfo(valueSpec *ast.ValueSpec, pkg *packages.Package) *ConstInfo {
	return &ConstInfo{
		Name:  valueSpec.Names[0].Name,
		Value: pkg.TypesInfo.ObjectOf(valueSpec.Names[0]).(*types.Const).Val().ExactString(),
		Doc:   valueSpec.Doc.Text() + valueSpec.Comment.Text(),
	}
}

// String implements the stringer interface.
func (c *ConstInfo) String() string {
	return fmt.Sprintf("Name: \"%s\"\nValue: %s\nDoc: \"%s\"\n",
		c.Name, c.Value, strings.ReplaceAll(c.Doc, "\n", "\\n"))
}

// InlineDoc replaces all new line in the documentation with spaces, making the resulting string
// suitable to be placed inline.
func (c *ConstInfo) InlineDoc() string {
	return strings.TrimRight(strings.ReplaceAll(c.Doc, "\n", " "), " ")
}
