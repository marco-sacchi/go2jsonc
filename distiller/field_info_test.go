package distiller

import (
	"fmt"
	"github.com/marco-sacchi/go2jsonc/testutils"
	"go/ast"
	"reflect"
	"strings"
	"testing"
)

type FieldInfoMatch struct {
	Type       string
	Name       string
	Layout     FieldLayout
	IsEmbedded bool
	Tags       map[string]string
	Doc        string
}

func (f *FieldInfoMatch) String() string {
	return fmt.Sprintf("Type: %s\nName: \"%s\"\nLayout: %v\nIsEmbedded: %v\nTags: %+v\nDoc: \"%v\"\n",
		f.Type, f.Name, f.Layout, f.IsEmbedded, f.Tags, strings.ReplaceAll(f.Doc, "\n", "\\n"))
}

func TestFieldInfo_String(t *testing.T) {
	info := getFieldsInfo(t, []string{"../testdata"})
	want := []string{
		`Type: int
Name: "Identifier"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:id]
Doc: "Identifier documentation block.\n"
`,
		`Type: bool
Name: "Enabled"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Enabled comment line.\n"
`,
		`Type: uint32
Name: "Reserved"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:reserved]
Doc: ""
`,
		`Type: github.com/marco-sacchi/go2jsonc/testdata.Embedded
Name: ""
Layout: 0
Element type: <nil>
IsEmbedded: true
Tags: map[]
Doc: "Embedded documentation block.\n"
`,
		`Type: float32
Name: "Position"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:position]
Doc: "Position comment line.\n"
`,
		`Type: float32
Name: "Velocity"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:velocity]
Doc: "Velocity documentation block.\n"
`,
		`Type: float32
Name: "Acceleration"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:accel]
Doc: ""
`,
		`Type: string
Name: "Reserved"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:reserved]
Doc: "Shadowing field.\n"
`,
		`Type: string
Name: "Name"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Name describes the protocol name.\nMultiple line documentation test.\nProtocol name.\n"
`,
		`Type: int
Name: "Major"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Major version.\n"
`,
		`Type: int
Name: "Minor"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Minor version.\n"
`,
		`Type: string
Name: "IP"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Remote IP address.\n"
`,
		`Type: int
Name: "Port"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Remote port.\n"
`,
		`Type: github.com/marco-sacchi/go2jsonc/testdata.Protocol
Name: "Default"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:default_proto]
Doc: "Default protocol.\n"
`,
		`Type: []github.com/marco-sacchi/go2jsonc/testdata.Protocol
Name: "Optionals"
Layout: 1
Element type: github.com/marco-sacchi/go2jsonc/testdata.Protocol
IsEmbedded: false
Tags: map[json:optional_protos]
Doc: "Optional supported protocols.\n"
`,
		`Type: string
Name: "Name"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "User name documentation block.\nUser name comment.\n"
`,
		`Type: string
Name: "Surname"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "User surname comment.\n"
`,
		`Type: int
Name: "Age"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:age]
Doc: "Age documentation block.\nUser age.\n"
`,
		`Type: int
Name: "StarsCount"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[json:stars_count]
Doc: "Number of stars achieved.\n"
`,
		`Type: []string
Name: "Addresses"
Layout: 1
Element type: string
IsEmbedded: false
Tags: map[]
Doc: "Addresses comment.\n"
`,
		`Type: map[string]string
Name: "Tags"
Layout: 2
Element type: string
IsEmbedded: false
Tags: map[]
Doc: "User tags.\n"
`,
		`Type: github.com/marco-sacchi/go2jsonc/testdata.ConstType
Name: "Type"
Layout: 0
Element type: <nil>
IsEmbedded: false
Tags: map[]
Doc: "Type documentation block.\nType of constant.\n"
`,
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for i, fieldInfo := range info {
		s := fieldInfo.String()
		if s != want[i] {
			t.Fatalf("String return mismatch: got:\n%s\nwant:\n%s\n",
				whitespacesReplacer.Replace(s),
				whitespacesReplacer.Replace(want[i]))
		}
	}
}

func TestFieldInfo_FormatDoc(t *testing.T) {
	info := getFieldsInfo(t, []string{"../testdata"})
	testTable := []struct {
		types bool
		want  string
	}{
		{types: true, want: "// int - Identifier documentation block.\n"},
		{types: true, want: "// bool - Enabled comment line.\n"},
		{types: true, want: "// uint32\n"},
		{types: true, want: "// testdata.Embedded - Embedded documentation block.\n"},
		{types: true, want: "// float32 - Position comment line.\n"},
		{types: true, want: "// float32 - Velocity documentation block.\n"},
		{types: true, want: "// float32\n"},
		{types: true, want: "// string - Shadowing field.\n"},
		{types: true, want: "// string - Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{types: true, want: "// int - Major version.\n"},
		{types: true, want: "// int - Minor version.\n"},
		{types: true, want: "// string - Remote IP address.\n"},
		{types: true, want: "// int - Remote port.\n"},
		{types: true, want: "// testdata.Protocol - Default protocol.\n"},
		{types: true, want: "// []testdata.Protocol - Optional supported protocols.\n"},
		{types: true, want: "// string - User name documentation block.\n// User name comment.\n"},
		{types: true, want: "// string - User surname comment.\n"},
		{types: true, want: "// int - Age documentation block.\n// User age.\n"},
		{types: true, want: "// int - Number of stars achieved.\n"},
		{types: true, want: "// []string - Addresses comment.\n"},
		{types: true, want: "// map[string]string - User tags.\n"},
		{
			types: true,
			want: `// testdata.ConstType - Type documentation block.
// Type of constant.
// Allowed values:
// ConstTypeA =   0  ConstTypeA doc block. ConstTypeA comment.
// ConstTypeB =   1  ConstTypeB comment.
// ConstTypeC =   2  ConstTypeC doc block. ConstTypeC comment.
// ConstTypeD =  32  ConstTypeD doc block.
// ConstTypeE =  64  ConstTypeE doc block. ConstTypeE comment.
// ConstTypeF = 128  ConstTypeF doc block. ConstTypeF comment.
`,
		},

		{types: false, want: "// Identifier documentation block.\n"},
		{types: false, want: "// Enabled comment line.\n"},
		{types: false, want: ""},
		{types: false, want: "// Embedded documentation block.\n"},
		{types: false, want: "// Position comment line.\n"},
		{types: false, want: "// Velocity documentation block.\n"},
		{types: false, want: ""},
		{types: false, want: "// Shadowing field.\n"},
		{types: false, want: "// Name describes the protocol name.\n// Multiple line documentation test.\n// Protocol name.\n"},
		{types: false, want: "// Major version.\n"},
		{types: false, want: "// Minor version.\n"},
		{types: false, want: "// Remote IP address.\n"},
		{types: false, want: "// Remote port.\n"},
		{types: false, want: "// Default protocol.\n"},
		{types: false, want: "// Optional supported protocols.\n"},
		{types: false, want: "// User name documentation block.\n// User name comment.\n"},
		{types: false, want: "// User surname comment.\n"},
		{types: false, want: "// Age documentation block.\n// User age.\n"},
		{types: false, want: "// Number of stars achieved.\n"},
		{types: false, want: "// Addresses comment.\n"},
		{types: false, want: "// User tags.\n"},
		{
			types: false,
			want: `// Type documentation block.
// Type of constant.
// Allowed values:
// ConstTypeA =   0  ConstTypeA doc block. ConstTypeA comment.
// ConstTypeB =   1  ConstTypeB comment.
// ConstTypeC =   2  ConstTypeC doc block. ConstTypeC comment.
// ConstTypeD =  32  ConstTypeD doc block.
// ConstTypeE =  64  ConstTypeE doc block. ConstTypeE comment.
// ConstTypeF = 128  ConstTypeF doc block. ConstTypeF comment.
`,
		},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "———➞")
	for i, test := range testTable {
		doc := info[i%len(info)].FormatDoc("", test.types)
		if doc != test.want {
			t.Fatalf("FormatDoc return mismatch:\ngot:%v\nwant:%v\n",
				whitespacesReplacer.Replace(doc),
				whitespacesReplacer.Replace(test.want))
		}
	}
}

func TestFieldInfo(t *testing.T) {
	dirs := []string{"../testdata"}
	testFieldInfo(t, dirs, []*FieldInfoMatch{
		// testdata/embedding.go
		{
			Type: "int", Name: "Identifier",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "id"},
			Doc:  "Identifier documentation block.\n",
		},
		{
			Type: "bool", Name: "Enabled",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Enabled comment line.\n",
		},
		{
			Type: "uint32", Name: "Reserved",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.Embedded", Name: "",
			Layout: LayoutSingle, IsEmbedded: true,
			Tags: nil,
			Doc:  "Embedded documentation block.\n",
		},
		{
			Type: "float32", Name: "Position",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "position"},
			Doc:  "Position comment line.\n",
		},
		{
			Type: "float32", Name: "Velocity",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "velocity"},
			Doc:  "Velocity documentation block.\n",
		},
		{
			Type: "float32", Name: "Acceleration",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "accel"},
			Doc:  "",
		},
		{
			Type: "string", Name: "Reserved",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "Shadowing field.\n",
		},
		// testdata/empty.go
		// testdata/nesting.go
		{
			Type: "string", Name: "Name",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Name describes the protocol name.\nMultiple line documentation test.\nProtocol name.\n",
		},
		{
			Type: "int", Name: "Major",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Major version.\n",
		},
		{
			Type: "int", Name: "Minor",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Minor version.\n",
		},
		{
			Type: "string", Name: "IP",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote IP address.\n",
		},
		{
			Type: "int", Name: "Port",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote port.\n",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.Protocol", Name: "Default",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "default_proto"},
			Doc:  "Default protocol.\n",
		},
		{
			Type: "[]github.com/marco-sacchi/go2jsonc/testdata.Protocol", Name: "Optionals",
			Layout: LayoutArray, IsEmbedded: false,
			Tags: map[string]string{"json": "optional_protos"},
			Doc:  "Optional supported protocols.\n",
		},
		// testdata/simple.go
		{
			Type: "string", Name: "Name",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "User name documentation block.\nUser name comment.\n",
		},
		{
			Type: "string", Name: "Surname",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "User surname comment.\n",
		},
		{
			Type: "int", Name: "Age",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "age"},
			Doc:  "Age documentation block.\nUser age.\n",
		},
		{
			Type: "int", Name: "StarsCount",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "stars_count"},
			Doc:  "Number of stars achieved.\n",
		},
		{
			Type: "[]string", Name: "Addresses",
			Layout: LayoutArray, IsEmbedded: false,
			Tags: nil,
			Doc:  "Addresses comment.\n",
		},
		{
			Type: "map[string]string", Name: "Tags",
			Layout: LayoutMap, IsEmbedded: false,
			Tags: nil,
			Doc:  "User tags.\n",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.ConstType", Name: "Type",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Type documentation block.\nType of constant.\n",
		},
	})
}

func TestFieldInfoMultiPackage(t *testing.T) {
	dirs := []string{
		"../testdata/multipkg/network",
		"../testdata/multipkg/stats",
		"../testdata/multipkg",
	}
	testFieldInfo(t, dirs, []*FieldInfoMatch{
		// testdata/multipkg/status.go
		{
			Type: "bool", Name: "Connected",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connected flag comment.\n",
		},
		{
			Type:   "github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.ConnState",
			Name:   "State",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connection state comment.\n",
		},
		// testdata/multipkg/info.go
		{
			Type: "int", Name: "PacketLoss",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "packet_loss"},
			Doc:  "PacketLoss documentation block.\nPacket loss comment.\n",
		},
		{
			Type: "int", Name: "RoundTripTime",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: map[string]string{"json": "round_trip_time"},
			Doc:  "Round-trip time in milliseconds.\n",
		},
		// testdata/multipkg/multi_package.go
		{
			Type:   "github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.Status",
			Name:   "NetStatus",
			Layout: LayoutSingle, IsEmbedded: false,
			Tags: nil,
			Doc:  "Network status.\n",
		},
		{
			Type:   "github.com/marco-sacchi/go2jsonc/testdata/multipkg/stats.Info",
			Name:   "",
			Layout: LayoutSingle, IsEmbedded: true,
			Tags: nil,
			Doc:  "Statistics info.\n",
		},
	})
}

func testFieldInfo(t *testing.T, patterns []string, want []*FieldInfoMatch) {
	fields := getFieldsInfo(t, patterns)

	if len(fields) != len(want) {
		t.Fatalf("Parsed %d fields, want %d.", len(fields), len(want))
	}

	for i, field := range fields {
		if field.Type.String() != want[i].Type || field.Name != want[i].Name ||
			field.Layout != want[i].Layout || field.IsEmbedded != want[i].IsEmbedded ||
			!reflect.DeepEqual(field.Tags, want[i].Tags) || field.Doc != want[i].Doc {
			t.Fatalf("Parsed field mismatch:\n%s\n\nwant:\n%s\n", field, want[i])
		}
	}
}

func getFieldsInfo(t *testing.T, patterns []string) []*FieldInfo {
	pkgs := testutils.LoadPackage(t, patterns...)
	var fields []*FieldInfo

	for _, pattern := range patterns {
		_, err := NewPackageInfo(pattern, "")
		if err != nil {
			t.Fatalf("Error loading package %s: %v", pattern, err)
		}
	}

	for _, pkg := range pkgs {
		for _, astFile := range pkg.Syntax {
			for _, decl := range astFile.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				ast.Inspect(genDecl, func(node ast.Node) bool {
					var field *ast.Field
					field, ok = node.(*ast.Field)
					if !ok {
						return true
					}

					fields = append(fields, NewFieldInfo(field, pkg))

					return true
				})
			}
		}
	}
	return fields
}
