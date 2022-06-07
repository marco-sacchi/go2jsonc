package testutils

import (
	"fmt"
	"golang.org/x/tools/go/packages"
	"strings"
	"testing"
)

type FieldInfo struct {
	Type       string
	Name       string
	IsArray    bool
	IsEmbedded bool
	Tags       map[string]string
	Doc        string
}

func (f *FieldInfo) String() string {
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nIsArray: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type, f.Name, f.IsArray, f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}

type StructInfo struct {
	Package     string
	Name        string
	FieldsCount int
	Doc         string
	Defaults    map[string]interface{}
}

func (s *StructInfo) String() string {
	return fmt.Sprintf("Package: %s\nName: \"%s\"\nFieldCount: %v\nDoc: \"%v\"\nDefaults: %+v",
		s.Package, s.Name, s.FieldsCount, strings.ReplaceAll(s.Doc, "\n", "\\n"), s.Defaults)
}

func LoadPackage(t *testing.T, patterns ...string) []*packages.Package {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) == 0 {
		t.Fatalf("Expected at least one package reading: %v", patterns)
	}

	return pkgs
}
