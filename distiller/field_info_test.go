package distiller

import (
	"github.com/marco-sacchi/go2jsonc/testutils"
	"go/ast"
	"reflect"
	"testing"
)

func TestFieldInfo(t *testing.T) {
	dirs := []string{"../testdata"}
	testFieldInfo(t, dirs, []*testutils.FieldInfo{
		// testdata/embedding.go
		{
			Type: "int", Name: "Identifier",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "id"},
			Doc:  "Identifier documentation block.\n",
		},
		{
			Type: "bool", Name: "Enabled",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Enabled comment line.\n",
		},
		{
			Type: "uint32", Name: "Reserved",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.Embedded", Name: "",
			IsArray: false, IsEmbedded: true,
			Tags: nil,
			Doc:  "Embedded documentation block.\n",
		},
		{
			Type: "float32", Name: "Position",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "position"},
			Doc:  "Position comment line.\n",
		},
		{
			Type: "float32", Name: "Velocity",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "velocity"},
			Doc:  "Velocity documentation block.\n",
		},
		{
			Type: "float32", Name: "Acceleration",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "accel"},
			Doc:  "",
		},
		{
			Type: "string", Name: "Reserved",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "reserved"},
			Doc:  "Shadowing field.\n",
		},
		// testdata/empty.go
		// testdata/nesting.go
		{
			Type: "string", Name: "Name",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Name describes the protocol name.\nMultiple line documentation test.\nProtocol name.\n",
		},
		{
			Type: "int", Name: "Major",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Major version.\n",
		},
		{
			Type: "int", Name: "Minor",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Minor version.\n",
		},
		{
			Type: "string", Name: "IP",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote IP address.\n",
		},
		{
			Type: "int", Name: "Port",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Remote port.\n",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.Protocol", Name: "Default",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "default_proto"},
			Doc:  "Default protocol.\n",
		},
		{
			Type: "github.com/marco-sacchi/go2jsonc/testdata.Protocol", Name: "Optionals",
			IsArray: true, IsEmbedded: false,
			Tags: map[string]string{"json": "optional_protos"},
			Doc:  "Optional supported protocols.\n",
		},
		// testdata/simple.go
		{
			Type: "string", Name: "Name",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "User name documentation block.\nUser name comment.\n",
		},
		{
			Type: "string", Name: "Surname",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "User surname comment.\n",
		},
		{
			Type: "int", Name: "Age",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "age"},
			Doc:  "Age documentation block.\nUser age.\n",
		},
		{
			Type: "int", Name: "StarsCount",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "stars_count"},
			Doc:  "Number of stars achieved.\n",
		},
		{
			Type: "string", Name: "Addresses",
			IsArray: true, IsEmbedded: false,
			Tags: nil,
			Doc:  "Addresses comment.\n",
		},
	})
}

func TestFieldInfoMultiPackage(t *testing.T) {
	dirs := []string{
		"../testdata/multipkg/network",
		"../testdata/multipkg/stats",
		"../testdata/multipkg",
	}
	testFieldInfo(t, dirs, []*testutils.FieldInfo{
		// testdata/multipkg/status.go
		{
			Type: "bool", Name: "Connected",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connected flag comment.\n",
		},
		{
			Type:    "github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.ConnState",
			Name:    "State",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Connection state comment.\n",
		},
		// testdata/multipkg/info.go
		{
			Type: "int", Name: "PacketLoss",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "packet_loss"},
			Doc:  "PacketLoss documentation block.\nPacket loss comment.\n",
		},
		{
			Type: "int", Name: "RoundTripTime",
			IsArray: false, IsEmbedded: false,
			Tags: map[string]string{"json": "round_trip_time"},
			Doc:  "Round-trip time in milliseconds.\n",
		},
		// testdata/multipkg/multi_package.go
		{
			Type:    "github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.Status",
			Name:    "NetStatus",
			IsArray: false, IsEmbedded: false,
			Tags: nil,
			Doc:  "Network status.\n",
		},
		{
			Type:    "github.com/marco-sacchi/go2jsonc/testdata/multipkg/stats.Info",
			Name:    "",
			IsArray: false, IsEmbedded: true,
			Tags: nil,
			Doc:  "Statistics info.\n",
		},
	})
}

func testFieldInfo(t *testing.T, patterns []string, want []*testutils.FieldInfo) {
	pkgs := testutils.LoadPackage(t, patterns...)

	var fields []*FieldInfo

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

	if len(fields) != len(want) {
		t.Fatalf("Parsed %d fields, want %d.", len(fields), len(want))
	}

	for i, field := range fields {
		if field.Type.String() != want[i].Type || field.Name != want[i].Name ||
			field.IsArray != want[i].IsArray || field.IsEmbedded != want[i].IsEmbedded ||
			!reflect.DeepEqual(field.Tags, want[i].Tags) || field.Doc != want[i].Doc {
			t.Fatalf("Parsed field mismatch:\n%s\n\nwant:\n%s\n", field, want[i])
		}
	}
}
