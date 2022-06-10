package testutils

import (
	"golang.org/x/tools/go/packages"
	"testing"
)

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
