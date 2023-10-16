package go2jsonc

import (
	"github.com/marco-sacchi/go2jsonc/distiller"
	"go/constant"
	"go/types"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	var tests = []struct {
		pkgDir   string
		typeName string
		filename string
		mode     DocTypesMode
	}{
		{"./testdata", "Embedding", "./testdata/embedding.jsonc", AllFields},
		{"./testdata", "Empty", "./testdata/empty.jsonc", AllFields},
		{"./testdata", "Nesting", "./testdata/nesting.jsonc", AllFields},
		{"./testdata", "Simple", "./testdata/simple.jsonc", AllFields},
		{"./testdata/multipkg", "MultiPackage", "./testdata/multipkg/multi_package.jsonc", AllFields},

		{"./testdata", "Embedding", "./testdata/embedding_not_fields.jsonc", NotFields},
		{"./testdata", "Nesting", "./testdata/nesting_not_fields.jsonc", NotFields},
		{"./testdata", "Simple", "./testdata/simple_not_fields.jsonc", NotFields},
		{"./testdata/multipkg", "MultiPackage", "./testdata/multipkg/multi_package_not_fields.jsonc", NotFields},

		{"./testdata", "Embedding", "./testdata/embedding_not_struct.jsonc", NotStructFields},
		{"./testdata", "Nesting", "./testdata/nesting_not_struct.jsonc", NotStructFields},
		{"./testdata", "Simple", "./testdata/simple_not_struct.jsonc", NotStructFields},
		{"./testdata/multipkg", "MultiPackage", "./testdata/multipkg/multi_package_not_struct.jsonc", NotStructFields},

		{"./testdata", "Embedding", "./testdata/embedding_not_array.jsonc", NotArrayFields},
		{"./testdata", "Nesting", "./testdata/nesting_not_array.jsonc", NotArrayFields},
		{"./testdata", "Simple", "./testdata/simple_not_array.jsonc", NotArrayFields},
		{"./testdata/multipkg", "MultiPackage", "./testdata/multipkg/multi_package_not_array.jsonc", NotArrayFields},

		{"./testdata", "Embedding", "./testdata/embedding_not_map.jsonc", NotMapFields},
		{"./testdata", "Nesting", "./testdata/nesting_not_map.jsonc", NotMapFields},
		{"./testdata", "Simple", "./testdata/simple_not_map.jsonc", NotMapFields},
		{"./testdata/multipkg", "MultiPackage", "./testdata/multipkg/multi_package_not_map.jsonc", NotMapFields},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for _, test := range tests {
		jsonc, err := Generate(test.pkgDir, test.typeName, test.mode)
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(test.filename)
		if err != nil {
			t.Fatal(err)
		}

		want := string(content)

		if jsonc != want {
			t.Fatalf("Generated JSONC mismatch for %s struct:\n%s\n\nwant %s:\n%s",
				test.typeName,
				whitespacesReplacer.Replace(jsonc),
				test.filename,
				whitespacesReplacer.Replace(want))
		}
	}

	_, err := Generate("./testdata/invalid-path", "", AllFields)
	if err == nil {
		t.Fatalf("Generating for invalid path: expected error, got nil.")
	}

	_, err = Generate("./testdata", "invalid-struct", AllFields)
	if err == nil {
		t.Fatalf("Generating for invalid struct: expected error, got nil.")
	}
}

func TestGenerator_typeZero(t *testing.T) {
	tests := []struct {
		info *distiller.FieldInfo
		want interface{}
	}{
		{info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutArray}, want: make([]interface{}, 0)},
		{info: &distiller.FieldInfo{Type: nil, Layout: distiller.LayoutMap}, want: make(map[interface{}]interface{})},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Bool], Layout: distiller.LayoutSingle}, want: false},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int], Layout: distiller.LayoutSingle}, want: 0},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int8], Layout: distiller.LayoutSingle}, want: int8(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int16], Layout: distiller.LayoutSingle}, want: int16(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int32], Layout: distiller.LayoutSingle}, want: int32(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Int64], Layout: distiller.LayoutSingle}, want: int64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint], Layout: distiller.LayoutSingle}, want: uint(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint8], Layout: distiller.LayoutSingle}, want: uint8(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint16], Layout: distiller.LayoutSingle}, want: uint16(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint32], Layout: distiller.LayoutSingle}, want: uint32(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uint64], Layout: distiller.LayoutSingle}, want: uint64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Uintptr], Layout: distiller.LayoutSingle}, want: uintptr(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Float32], Layout: distiller.LayoutSingle}, want: float32(0.0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Float64], Layout: distiller.LayoutSingle}, want: 0.0},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Complex64], Layout: distiller.LayoutSingle}, want: complex64(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.Complex128], Layout: distiller.LayoutSingle}, want: complex128(0)},
		{info: &distiller.FieldInfo{Type: types.Typ[types.String], Layout: distiller.LayoutSingle}, want: constant.MakeString("")},
	}

	for _, test := range tests {
		zero := typeZero(test.info)

		if !reflect.DeepEqual(zero, test.want) {
			t.Fatalf("Zero value mismatch for type %v: got %v, want %v",
				test.info.Type.String(), zero, test.want)
		}
	}
}
