package main

import (
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"mmatt/elicit"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

type stepImpl struct {
	pattern *regexp.Regexp
	fn      string
}

func parseStepPackage(path string) []stepImpl {

	pkg, patterns := parseStepImplPackageDir(path)

	stepsFound := make([]stepImpl, 0)

	for _, n := range pkg.Scope().Names() {

		obj := pkg.Scope().Lookup(n)

		if !isStepImpl(obj) {
			continue
		}

		p, hasPattern := patterns[n]
		if !hasPattern {
			continue
		}

		r, err := regexp.Compile(p)

		if err != nil {
			log.Fatalf("parsing step expression for %s - %s: %s", n, p, err)
		}

		stepsFound = append(stepsFound, stepImpl{
			pattern: r,
			fn:      n,
		})
	}

	return stepsFound
}

func isStepImpl(obj types.Object) bool {
	typ := obj.Type()
	if typ == nil {
		return false
	}

	if _, isFunc := obj.(*types.Func); !isFunc {
		return false
	}

	sig, isSig := typ.(*types.Signature)
	if !isSig {
		return false
	}

	if sig.Recv() != nil ||
		sig.Results() != nil ||
		sig.Params().Len() == 0 {
		return false
	}

	if !isFirstParamElicitContext(sig) {
		return false
	}

	return true
}

func isFirstParamElicitContext(sig *types.Signature) bool {
	p0Typ := sig.Params().At(0).Type()

	pointer, isPtr := p0Typ.(*types.Pointer)

	if !isPtr {
		return false
	}

	named, isNamed := pointer.Elem().(*types.Named)

	if !isNamed {
		return false
	}

	// TODO: There must be a better way of checking the formal parameter type against the expected one.
	pkgPath := named.Obj().Pkg().Path()
	typeName := named.Obj().Name()

	elicitContextType := reflect.TypeOf((*elicit.Context)(nil)).Elem()
	if elicitContextType.PkgPath() != pkgPath || elicitContextType.Name() != typeName {
		return false
	}

	return true
}

func parseStepImplPackageDir(directory string) (*types.Package, map[string]string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("parsing step directory %q: %s", directory, err)
	}

	var names []string
	names = append(names, pkg.GoFiles...)
	names = prefixDirectory(directory, names)
	return parsePackage(directory, names)
}

func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

func parsePackage(directory string, names []string) (*types.Package, map[string]string) {
	var astFiles []*ast.File
	fs := token.NewFileSet()
	patterns := make(map[string]string)

	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, nil, parser.ParseComments)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}

		for k, v := range extractPatterns(parsedFile) {
			patterns[k] = v
		}

		astFiles = append(astFiles, parsedFile)
	}

	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}

	return typeCheckPkg(directory, fs, astFiles), patterns
}

func extractPatterns(astFile *ast.File) map[string]string {
	patterns := make(map[string]string)

	for _, d := range astFile.Decls {
		if f, isFuncDecl := d.(*ast.FuncDecl); isFuncDecl {
			for _, t := range f.Doc.List {
				if strings.HasPrefix(t.Text, "//elicit:step") {
					patterns[f.Name.Name] = strings.TrimSpace(strings.TrimPrefix(t.Text, "//elicit:step"))
				}
			}
		}
	}

	return patterns
}

func typeCheckPkg(dir string, fs *token.FileSet, astFiles []*ast.File) *types.Package {
	defs := make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: importer.Default(), FakeImportC: true}
	info := &types.Info{
		Defs: defs,
	}
	typesPkg, err := config.Check(dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	return typesPkg
}
