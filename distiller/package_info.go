package distiller

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/packages"
	"os"
	"path/filepath"
)

// PackageInfo holds information about a package.
type PackageInfo struct {
	Package     *packages.Package
	Imported    []*packages.Package     // Packages imported from this.
	Structs     map[string]*StructInfo  // Structures declared on this package mapped by fully qualified name.
	TypedConsts map[string][]*ConstInfo // Typed constants grouped by fully qualified type name.
}

// loadedPackages caches the loaded/imported packages.
var loadedPackages = make(map[string]*PackageInfo)

// LookupStruct searches loaded packages for the specified fully qualified struct name.
// It returns nil in case of no matches.
func LookupStruct(name string) *StructInfo {
	for _, pkg := range loadedPackages {
		s, ok := pkg.Structs[name]
		if ok {
			return s
		}
	}

	return nil
}

// LookupTypedConsts searches loaded packages for declared constants of specified fully qualified named type.
// It returns nil in case of no matches.
func LookupTypedConsts(name string) []*ConstInfo {
	consts := []*ConstInfo(nil)
	for _, pkg := range loadedPackages {
		c, ok := pkg.TypedConsts[name]
		if ok {
			consts = append(consts, c...)
		}
	}

	return consts
}

// NewPackageInfo creates a new package information object from given directory. The passed name defines
// the struct for which read also defaults values.
func NewPackageInfo(dir string, typeName string) (*PackageInfo, error) {
	pkgInfo := new(PackageInfo)
	pkgInfo.Structs = make(map[string]*StructInfo)
	pkgInfo.TypedConsts = make(map[string][]*ConstInfo)

	err := pkgInfo.readPackage(dir)
	if err != nil {
		return nil, err
	}

	if typeName != "" {
		s, ok := pkgInfo.Structs[pkgInfo.Package.PkgPath+"."+typeName]
		if !ok {
			return nil, fmt.Errorf("cannot find struct %s in package %s", typeName, pkgInfo.Package.PkgPath)
		}

		if err = s.ParseDefaultsMethod(); err != nil {
			return nil, err
		}
	}

	loadedPackages[pkgInfo.Package.PkgPath] = pkgInfo

	return pkgInfo, nil
}

// readPackage reads information for the package defined in the given directory and all imported packages.
func (p *PackageInfo) readPackage(dir string) error {
	ok, err := isDirectory(dir)
	if err != nil {
		return err
	}

	if !ok {
		return fmt.Errorf("%v is not a directory", dir)
	}

	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedFiles |
			packages.NeedImports | packages.NeedDeps | packages.NeedSyntax |
			packages.NeedTypes | packages.NeedTypesInfo,
	}

	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return err
	}

	if len(pkgs) != 1 {
		return fmt.Errorf("expected exactly 1 package for given pattern")
	}

	p.Package = pkgs[0]

	for ident, object := range p.Package.TypesInfo.Defs {
		typeName, ok := object.(*types.TypeName)
		if !ok {
			continue
		}

		for _, astFile := range p.Package.Syntax {
			// Use the position to test if the type is declared in this file.
			if astFile.Pos() >= typeName.Pos() || typeName.Pos() >= astFile.End() {
				continue
			}

			typeNameString := typeName.Type().String()

			nodes, _ := astutil.PathEnclosingInterval(astFile, typeName.Pos(), typeName.Pos())
			isStruct := false
			for _, node := range nodes {
				var genDecl *ast.GenDecl
				genDecl, ok = node.(*ast.GenDecl)

				if !ok {
					continue
				}

				if len(genDecl.Specs) != 1 {
					return fmt.Errorf("expected one spec for struct declaration")
				}

				var typeSpec *ast.TypeSpec
				typeSpec, ok = genDecl.Specs[0].(*ast.TypeSpec)
				// Identifier is not a type declaration or not match, continue.
				if !ok || typeSpec.Name != ident {
					continue
				}

				// The identifier matches but the type is not a struct, exit loop.
				if _, ok = typeSpec.Type.(*ast.StructType); !ok {
					break
				}

				info := NewStructInfo(genDecl, p.Package)
				for _, field := range info.Fields {
					var namedType *types.Named
					namedType, ok = field.Type.(*types.Named)

					// types.Basic type.
					if !ok {
						continue
					}

					// Check if required package is loaded.
					pkgPath := namedType.Obj().Pkg().Path()
					_, ok = loadedPackages[pkgPath]
					if pkgPath == p.Package.PkgPath || ok {
						continue
					}

					// Load required package.
					imported := p.Package.Imports[pkgPath]
					if _, err = NewPackageInfo(filepath.Dir(imported.GoFiles[0]), ""); err != nil {
						return err
					}
				}

				p.Structs[typeNameString] = info
				isStruct = true
				break
			}

			// Not a struct, check if identifier is used by typed constants.
			if !isStruct {
				if consts := p.readIdentConsts(astFile, ident); consts != nil {
					p.TypedConsts[typeNameString] = consts
				}
			}
		}
	}

	return nil
}

// readIdentConsts reads typed constants what uses ident type. Returns nil if no constants use the type.
func (p *PackageInfo) readIdentConsts(astFile *ast.File, ident *ast.Ident) []*ConstInfo {
	consts := []*ConstInfo(nil)
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

			if valueSpec.Type != nil {
				var constIdent *ast.Ident
				constIdent, ok = valueSpec.Type.(*ast.Ident)
				if !ok || constIdent.Obj != ident.Obj {
					continue
				}
			}

			consts = append(consts, NewConstInfo(valueSpec, p.Package))
		}
	}

	return consts
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) (bool, error) {
	info, err := os.Stat(name)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
