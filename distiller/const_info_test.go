package distiller

import (
	"github.com/marco-sacchi/go2jsonc/testutils"
	"go/ast"
	"testing"
)

func TestConstInfo(t *testing.T) {
	pkgs := testutils.LoadPackage(t, "../testdata/consts.go")

	var consts []*ConstInfo

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range genDecl.Specs {
					var valueSpec *ast.ValueSpec
					valueSpec, ok = spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					consts = append(consts, NewConstInfo(valueSpec, pkg))
				}
			}
		}
	}

	want := []*ConstInfo{
		{Name: "ConstTypeA", Value: "0", Doc: "ConstTypeA doc block.\nConstTypeA comment.\n"},
		{Name: "ConstTypeB", Value: "1", Doc: "ConstTypeB comment.\n"},
		{Name: "ConstTypeC", Value: "2", Doc: "ConstTypeC doc block.\nConstTypeC comment.\n"},
		{Name: "ConstTypeD", Value: "32", Doc: "ConstTypeD doc block.\n"},
		{Name: "ConstTypeE", Value: "64", Doc: "ConstTypeE doc block.\nConstTypeE comment.\n"},
		{Name: "ConstTypeF", Value: "128", Doc: "ConstTypeF doc block.\nConstTypeF comment.\n"},
	}

	if len(consts) != len(want) {
		t.Fatalf("Parsed %d constants, want %d.", len(consts), len(want))
	}

	for i, constInfo := range consts {
		if constInfo.Name != want[i].Name || constInfo.Value != want[i].Value || constInfo.Doc != want[i].Doc {
			t.Fatalf("Parsed const mismatch:\n%s\n\nwant:\n%s\n", constInfo, want[i])
		}
	}
}
