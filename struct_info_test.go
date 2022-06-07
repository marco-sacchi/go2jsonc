package main

import (
	"github.com/marco-sacchi/go2jsonc/testdata/multipkg/network"
	"github.com/marco-sacchi/go2jsonc/testutils"
	"go/ast"
	"go/constant"
	"reflect"
	"testing"
)

func TestStructInfo(t *testing.T) {
	testStructInfo(t, "./testdata", []*testutils.StructInfo{
		// testdata/embedding.go
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Embedded",
			Doc:         "Embedded test struct.\n",
			FieldsCount: 3,
			Defaults:    nil,
		},
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Embedding",
			Doc:         "Embedding test struct.\n",
			FieldsCount: 5,
			Defaults:    nil,
		},
		// testdata/empty.go
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Empty",
			Doc:         "Empty empty test struct.\n",
			FieldsCount: 0,
			Defaults:    nil,
		},
		// testdata/nesting.go
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Protocol",
			Doc:         "Protocol defines a network protocol and version.\n",
			FieldsCount: 3,
			Defaults:    nil,
		},
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Nesting",
			Doc:         "Nesting checks for correct struct nesting.\n",
			FieldsCount: 4,
			Defaults:    nil,
		},
		// testdata/simple.go
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata",
			Name:        "Simple",
			Doc:         "Simple defines a simple user.\n",
			FieldsCount: 5,
			Defaults:    nil,
		},
	})
}

func TestStructInfoMultiPackage(t *testing.T) {
	testStructInfo(t, "./testdata/multipkg", []*testutils.StructInfo{
		// testdata/multipkg/multi_package.go
		{
			Package:     "github.com/marco-sacchi/go2jsonc/testdata/multipkg",
			Name:        "MultiPackage",
			Doc:         "MultiPackage tests the multi-package and import aliasing case.\n",
			FieldsCount: 2,
			Defaults:    nil,
		},
	})
}

func TestStructInfoDefaults(t *testing.T) {
	testStructInfoDefaults(t, "./testdata", "Simple", map[string]interface{}{
		"Name":       constant.MakeString("John"),
		"Surname":    constant.MakeString("Doe"),
		"Age":        constant.MakeInt64(30),
		"StarsCount": constant.MakeInt64(5),
		"Addresses": []interface{}{
			constant.MakeString("Address 1"),
			constant.MakeString("Address 2"),
			constant.MakeString("Address 3"),
		},
	})
}

func TestStructInfoDefaultsMultiPackage(t *testing.T) {
	testStructInfoDefaults(t, "./testdata/multipkg", "MultiPackage", map[string]interface{}{
		"NetStatus": map[string]interface{}{
			"Connected": constant.MakeBool(true),
			"State":     constant.MakeInt64(int64(network.StateDisconnected)),
		},
		"Info": map[string]interface{}{
			"PacketLoss":    constant.MakeInt64(64),
			"RoundTripTime": constant.MakeInt64(123),
		},
	})
}

func testStructInfo(t *testing.T, pattern string, want []*testutils.StructInfo) {
	pkgs := testutils.LoadPackage(t, pattern)

	var structs []*StructInfo

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var typeSpec *ast.TypeSpec
					typeSpec, ok = node.(*ast.TypeSpec)
					if !ok {
						return true
					}

					if _, ok = typeSpec.Type.(*ast.StructType); !ok {
						return true
					}

					structs = append(structs, NewStructInfo(genDecl, pkg))
					return true
				})
			}
		}
	}

	if len(structs) != len(want) {
		t.Fatalf("Parsed %d structs, want %d.", len(structs), len(want))
	}

	for i, s := range structs {
		if s.Package.PkgPath != want[i].Package ||
			s.Name != want[i].Name || s.Doc != want[i].Doc ||
			len(s.Fields) != want[i].FieldsCount ||
			!reflect.DeepEqual(s.Defaults, want[i].Defaults) {
			t.Fatalf("Parsed struct mismatch:\n%s\n\nwant:\n%s\n", s, want[i])
		}
	}
}

func testStructInfoDefaults(t *testing.T, pattern string, name string, want map[string]interface{}) {
	pkgs := testutils.LoadPackage(t, pattern)

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var typeSpec *ast.TypeSpec
					typeSpec, ok = node.(*ast.TypeSpec)
					if !ok {
						return true
					}

					if _, ok = typeSpec.Type.(*ast.StructType); !ok {
						return true
					}

					if typeSpec.Name.Name != name {
						return true
					}

					s := NewStructInfo(genDecl, pkg)
					if err := s.ParseDefaultsMethod(); err != nil {
						t.Fatal(err)
					}

					if !reflect.DeepEqual(s.Defaults, want) {
						t.Fatalf("Struct %s defaults mismatch:\n%+v\n\nwant:\n%+v", s.Name, s.Defaults, want)
					}

					return false
				})
			}
		}
	}
}
