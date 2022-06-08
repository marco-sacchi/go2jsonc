package distiller

import "testing"

func TestPackageInfo(t *testing.T) {
	info, err := NewPackageInfo("../testdata", "")
	if err != nil {
		t.Fatal(err)
	}

	wantConsts := []string{
		"github.com/marco-sacchi/go2jsonc/testdata.ConstType",
	}

	wantStructs := []string{
		"github.com/marco-sacchi/go2jsonc/testdata.Embedded",
		"github.com/marco-sacchi/go2jsonc/testdata.Embedding",
		"github.com/marco-sacchi/go2jsonc/testdata.Empty",
		"github.com/marco-sacchi/go2jsonc/testdata.Protocol",
		"github.com/marco-sacchi/go2jsonc/testdata.Nesting",
		"github.com/marco-sacchi/go2jsonc/testdata.Simple",
	}

	if len(info.TypedConsts) != len(wantConsts) {
		t.Fatalf("Parsed %d types used by typed constants, want %d",
			len(info.TypedConsts), len(wantConsts))
	}

	if len(info.Structs) != len(wantStructs) {
		t.Fatalf("Parsed %d structs, want %d", len(info.Structs), len(wantStructs))
	}

	for _, name := range wantConsts {
		if LookupTypedConsts(name) == nil {
			t.Fatalf("Cannot lookup typed constants of type %s", name)
		}
	}

	for _, name := range wantStructs {
		if LookupStruct(name) == nil {
			t.Fatalf("Cannot lookup struct %s", name)
		}
	}
}

func TestPackageInfoMultiPackage(t *testing.T) {
	_, err := NewPackageInfo("../testdata/multipkg", "")
	if err != nil {
		t.Fatal(err)
	}

	wantConsts := []string{
		"github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.ConnState",
	}

	wantStructs := []string{
		"github.com/marco-sacchi/go2jsonc/testdata/multipkg/network.Status",
		"github.com/marco-sacchi/go2jsonc/testdata/multipkg/stats.Info",
		"github.com/marco-sacchi/go2jsonc/testdata/multipkg.MultiPackage",
	}

	for _, name := range wantConsts {
		if LookupTypedConsts(name) == nil {
			t.Fatalf("Cannot lookup typed constants of type %s", name)
		}
	}

	for _, name := range wantStructs {
		if LookupStruct(name) == nil {
			t.Fatalf("Cannot lookup struct %s", name)
		}
	}
}
