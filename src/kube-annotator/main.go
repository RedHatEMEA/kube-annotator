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
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/types"
	"os"
	"reflect"
	"strings"
	_ "golang.org/x/tools/go/gcimporter"
)

// Don't recurse into structs of this type
var blacklist = map[string]bool{
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util.Time": true,
	"github.com/GoogleCloudPlatform/kubernetes/pkg/util.IntOrString": true,
}

func typefmt(typ types.Type) string {
	// Ugh.
	typename := typ.String()
	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		typename = strings.Replace(typename, p + "/src/", "", -1)
	}
	return typename
}

func print(s1 string, s2 string) {
	fmt.Printf("%-39s # %s\n", s1, s2)
}

func dump(typ types.Type, indent string) {
	if _, ok := typ.(*types.Named); ok {
		typ = typ.Underlying()
	}

	switch u := typ.(type) {
	case *types.Struct:
		for i := 0; i < u.NumFields(); i++ {
			st := reflect.StructTag(u.Tag(i))
			name := strings.Split(st.Get("json"), ",")[0]
			desc := st.Get("description")
			if desc != "" {
				desc = " (" + desc + ")"
			}
			
			if name == "" {
				dump(u.Field(i).Type(), indent)
			} else {
				print(indent + name + ":", typefmt(u.Field(i).Type()) + desc)
				if !blacklist[typefmt(u.Field(i).Type())] {
					dump(u.Field(i).Type(), indent + "  ")
				}
			}

			indent = strings.Replace(indent, "-", " ", -1)
		}

	case *types.Map:
		print(indent + "[" + typefmt(u.Key()) + "]:", typefmt(u.Elem()))

	case *types.Pointer:
		dump(u.Elem().Underlying(), indent)
		
	case *types.Slice:
		indent = strings.TrimSuffix(indent, "  ") + "- "
		if _, ok := u.Elem().Underlying().(*types.Struct); ok {
			dump(u.Elem().Underlying(), indent)
		} else {
			print(indent + "[" + typefmt(u.Elem()) + "]", "")
		}

	case *types.Basic:

	default:
		panic("unsupported")
	}
}

func importPkg(pkgname string) (*types.Package, error) {
	pkg, err := build.Import(pkgname, "", 0)
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	pkgmap, err := parser.ParseDir(fset, pkg.Dir, nil, 0)
	if err != nil {
		return nil, err
	}

	filelist := make([]*ast.File, 0, len(pkgmap[pkg.Name].Files))
	for _, f := range pkgmap[pkg.Name].Files {
		filelist = append(filelist, f)
	}

	config := types.Config{}
	return config.Check(pkg.Dir, fset, filelist, nil)
}

func walkPkg(pkg *types.Package) {
	for _, name := range pkg.Scope().Names() {
		obj := pkg.Scope().Lookup(name)

		if typename, ok := obj.(*types.TypeName); ok {
			named := typename.Type().(*types.Named)

			strukt, ok := named.Underlying().(*types.Struct)
			if !ok || strukt.NumFields() < 1 || strukt.Field(0).Name() != "TypeMeta" {
				continue
			}

			fmt.Printf("%s\n%s\n\n", typename.Name(), strings.Repeat("=", len(typename.Name())))
			dump(strukt, "")
			fmt.Printf("\n\n")
		}
	}
}

func main() {
	pkg, err := importPkg(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	
	walkPkg(pkg)
}
