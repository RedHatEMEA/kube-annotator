package main

/*
 * Copyright 2015 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/doc"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/types"
	"os"
	"reflect"
	"strings"
	_ "golang.org/x/tools/go/gcimporter"
)

type outputer interface {
	Struct(basename *types.TypeName, st reflect.StructTag, u *types.Struct, options string, tn string)
	Map(basename *types.TypeName, st reflect.StructTag, u *types.Map, options string, tn string)
	Slice(basename *types.TypeName, st reflect.StructTag, u *types.Slice, options string)
	Pointer(basename *types.TypeName, st reflect.StructTag, u *types.Pointer)
	Basic(basename *types.TypeName, st reflect.StructTag, u *types.Basic, options string, tn string)
}

func typefmt(typ types.Type) string {
	// Ugh.
	typename := typ.String()
	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		typename = strings.Replace(typename, p + "/src/", "", -1)
	}
	return typename
}

func dump(o outputer, basename *types.TypeName, typ types.Type, st reflect.StructTag) {
	tn := typefmt(typ)
	options := ""
	if named, ok := typ.(*types.Named); ok {
		tn = typefmt(named)
		options = getConsts(named)
		typ = typ.Underlying()
	}

	name := strings.Split(st.Get("json"), ",")[0]
	if _, ok := typ.(*types.Struct); !ok && name == "" {
		return
	}

	switch u := typ.(type) {
	case *types.Struct:
		o.Struct(basename, st, u, options, tn)

	case *types.Map:
		o.Map(basename, st, u, options, tn)

	case *types.Slice:
		o.Slice(basename, st, u, options)

	case *types.Pointer:
		o.Pointer(basename, st, u)

	case *types.Basic:
		o.Basic(basename, st, u, options, tn)

	default:
		panic("unsupported")
	}
}

func getConsts(named *types.Named) string {
	pkg := named.Obj().Pkg()

	s := make([]string, 0)
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)

		if konst, ok := obj.(*types.Const); ok {
			if konst.Type() == named {
				s = append(s, strings.Replace(konst.Val().String(), "\"", "", -1))
			}
		}
	}

	return strings.Join(s, "|")
}

func importPkg(pkgname string) (*types.Package, *ast.Package, error) {
	pkg, err := build.Import(pkgname, "", 0)
	if err != nil {
		return nil, nil, err
	}

	fset := token.NewFileSet()
	pkgmap, err := parser.ParseDir(fset, pkg.Dir, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	filelist := make([]*ast.File, 0, len(pkgmap[pkg.Name].Files))
	for _, f := range pkgmap[pkg.Name].Files {
		filelist = append(filelist, f)
	}

	config := types.Config{}
	typpkg, err := config.Check(pkg.Dir, fset, filelist, nil)

	return typpkg, pkgmap[pkg.Name], err
}

func walkPkg(typpkg *types.Package, docpkg *doc.Package) {
	for _, name := range typpkg.Scope().Names() {
		obj := typpkg.Scope().Lookup(name)

		if typename, ok := obj.(*types.TypeName); ok {
			named := typename.Type().(*types.Named)

			strukt, ok := named.Underlying().(*types.Struct)
			if !ok || strukt.NumFields() < 1 || strukt.Field(0).Name() != "TypeMeta" {
				continue
			}

			fmt.Printf("%s\n%s\n\n", typename.Name(), strings.Repeat("=", len(typename.Name())))

			for _, t := range docpkg.Types {
				if t.Name == typename.Name() {
					fmt.Printf("%v\n", t.Doc)
				}
			}

			o := DocOutput{}
			dump(&o, typename, strukt, reflect.StructTag("json:\"" + typename.Name() + "\""))
			fmt.Printf("\n\n")
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s package\n", os.Args[0])
		return
	}

	typpkg, astpkg, err := importPkg(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	docpkg := doc.New(astpkg, os.Args[1], 0)

	walkPkg(typpkg, docpkg)
}
