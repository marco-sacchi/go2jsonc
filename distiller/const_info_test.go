package distiller

import (
	"github.com/marco-sacchi/go2jsonc/testutils"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/packages"
	"strings"
	"testing"
)

func TestNewConstInfo(t *testing.T) {
	pkgs := testutils.LoadPackage(t, "../testdata/consts.go")
	consts := getConsts(pkgs)

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

func TestConstInfo_String(t *testing.T) {
	pkgs := testutils.LoadPackage(t, "../testdata/consts.go")
	consts := getConsts(pkgs)

	want := []string{
		"Name: \"ConstTypeA\"\nValue: 0\nDoc: \"ConstTypeA doc block.\\nConstTypeA comment.\\n\"\n",
		"Name: \"ConstTypeB\"\nValue: 1\nDoc: \"ConstTypeB comment.\\n\"\n",
		"Name: \"ConstTypeC\"\nValue: 2\nDoc: \"ConstTypeC doc block.\\nConstTypeC comment.\\n\"\n",
		"Name: \"ConstTypeD\"\nValue: 32\nDoc: \"ConstTypeD doc block.\\n\"\n",
		"Name: \"ConstTypeE\"\nValue: 64\nDoc: \"ConstTypeE doc block.\\nConstTypeE comment.\\n\"\n",
		"Name: \"ConstTypeF\"\nValue: 128\nDoc: \"ConstTypeF doc block.\\nConstTypeF comment.\\n\"\n",
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for i, constInfo := range consts {
		s := constInfo.String()
		if s != want[i] {
			t.Fatalf("String() return value mismatch: got:\n%s\nwant:\n%s\n",
				whitespacesReplacer.Replace(s),
				whitespacesReplacer.Replace(want[i]))
		}
	}
}

func TestConstInfo_InlineDoc(t *testing.T) {
	pkgs := testutils.LoadPackage(t, "../testdata/consts.go")
	consts := getConsts(pkgs)

	want := []string{
		"ConstTypeA doc block. ConstTypeA comment.",
		"ConstTypeB comment.",
		"ConstTypeC doc block. ConstTypeC comment.",
		"ConstTypeD doc block.",
		"ConstTypeE doc block. ConstTypeE comment.",
		"ConstTypeF doc block. ConstTypeF comment.",
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for i, constInfo := range consts {
		s := constInfo.InlineDoc()
		if s != want[i] {
			t.Fatalf("InlineDoc() return value mismatch: got:\n%s\nwant:\n%s\n",
				whitespacesReplacer.Replace(s),
				whitespacesReplacer.Replace(want[i]))
		}
	}
}

func getConsts(pkgs []*packages.Package) []*ConstInfo {
	var consts []*ConstInfo

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.CONST {
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
	return consts
}
