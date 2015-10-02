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
	_ "golang.org/x/tools/go/gcimporter"
	"golang.org/x/tools/go/types"
	"os"
)

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

	filelist := make([]*ast.File, 0)
	for _, f := range pkgmap[pkg.Name].Files {
		filelist = append(filelist, f)
	}

	config := types.Config{}
	typpkg, err := config.Check(pkg.Dir, fset, filelist, nil)

	return typpkg, pkgmap[pkg.Name], err
}

func walkPkg(typpkg *types.Package, docpkg *doc.Package, f func(*types.Struct, *types.TypeName, *doc.Package) string) {
	for _, name := range typpkg.Scope().Names() {
		obj := typpkg.Scope().Lookup(name)

		if typename, ok := obj.(*types.TypeName); ok {
			named := typename.Type().(*types.Named)

			if strukt, ok := named.Underlying().(*types.Struct); ok && strukt.NumFields() > 0 && strukt.Field(0).Name() == "TypeMeta" {
				if len(os.Args) == 3 || os.Args[3] == typename.Name() {
					fmt.Printf("%s\n", f(strukt, typename, docpkg))
				}
			}
		}
	}
}

func main() {
	if len(os.Args) < 3 || (os.Args[1] != "doc" && os.Args[1] != "jsonform") {
		fmt.Fprintf(os.Stderr, "usage: %s doc|jsonform package [type]\n", os.Args[0])
		return
	}

	typpkg, astpkg, err := importPkg(os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}

	docpkg := doc.New(astpkg, os.Args[1], 0)

	switch os.Args[1] {
	case "doc":
		walkPkg(typpkg, docpkg, docf)
	case "jsonform":
		walkPkg(typpkg, docpkg, jsonf)
	}
}
