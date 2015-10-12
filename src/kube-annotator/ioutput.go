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
	"os"
	"reflect"
	"sort"
	"strings"

	"golang.org/x/tools/go/types"
)

type IObj interface {
	Name() string
	Path() string
}

type IBase struct {
	name        string
	typ         string
	description string
	path        string
}

type IStruct struct {
	IBase
	items []IObj
}

type IMap struct {
	IBase
	keytyp string
	valtyp string
}

type ISlice struct {
	IBase
	items  []IObj
	valtyp string
}

type IBasic struct {
	IBase
	options string
}

func (o IStruct) Name() string { return o.name }
func (o IStruct) Path() string { return o.path }
func (o IMap) Name() string    { return o.name }
func (o IMap) Path() string    { return o.path }
func (o ISlice) Name() string  { return o.name }
func (o ISlice) Path() string  { return o.path }
func (o IBasic) Name() string  { return o.name }
func (o IBasic) Path() string  { return o.path }

func Struct(path string, st reflect.StructTag, u *types.Struct, named *types.Named) IObj {
	rv := IStruct{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = getname(named, u)
	rv.description = st.Get("description")
	rv.path = pathAppend(path, rv.name)

	for i := 0; i < u.NumFields(); i++ {
		rv2 := dump(rv.path, u.Field(i).Type(), reflect.StructTag(u.Tag(i)))
		if rv2 == nil {
			continue
		}
		if st, ok := rv2.(IStruct); ok && st.name == "" {
			for _, j := range st.items {
				rv.items = append(rv.items, j)
			}
		} else {
			rv.items = append(rv.items, rv2)
		}
	}

	return rv
}

func Map(path string, st reflect.StructTag, u *types.Map, named *types.Named) IObj {
	rv := IMap{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = getname(named, u)
	rv.description = st.Get("description")
	rv.path = pathAppend(path, rv.name)
	rv.keytyp = typefmt(u.Key())
	rv.valtyp = typefmt(u.Elem())

	return rv
}

func Slice(path string, st reflect.StructTag, u *types.Slice) IObj {
	und := u.Elem()
	if p, ok := und.(*types.Pointer); ok {
		und = p.Elem()
	}

	rv := ISlice{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = "[]" + typefmt(und)
	rv.description = st.Get("description")
	rv.path = pathAppend(path, rv.name)
	c := dump(rv.path, und, reflect.StructTag(""))
	if st, ok := c.(IStruct); ok {
		rv.items = st.items
	} else {
		rv.valtyp = typefmt(und)
	}

	return rv
}

func Pointer(path string, st reflect.StructTag, u *types.Pointer) IObj {
	return dump(path, u.Elem(), st)
}

func Basic(path string, st reflect.StructTag, u *types.Basic, named *types.Named) IObj {
	rv := IBasic{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = getname(named, u)
	rv.description = st.Get("description")
	rv.path = pathAppend(path, rv.name)
	if named != nil {
		rv.options = getConsts(named)
	}

	return rv
}

func getname(named *types.Named, typ types.Type) string {
	if named != nil {
		return typefmt(named)
	}
	return typefmt(typ)
}

func typefmt(typ types.Type) string {
	// Ugh.
	typename := typ.String()
	for _, p := range strings.Split(os.Getenv("GOPATH"), ":") {
		typename = strings.Replace(typename, p+"/src/", "", -1)
	}
	return typename
}

func dump(path string, typ types.Type, st reflect.StructTag) IObj {
	named, _ := typ.(*types.Named)
	if named != nil {
		typ = typ.Underlying()
	}

	if strings.Split(st.Get("json"), ",")[0] == "" {
		if _, ok := typ.(*types.Struct); !ok {
			if _, ok := typ.(*types.Pointer); !ok {
				return nil
			}
		}
	}

	switch u := typ.(type) {
	case *types.Struct:
		return Struct(path, st, u, named)

	case *types.Map:
		return Map(path, st, u, named)

	case *types.Slice:
		return Slice(path, st, u)

	case *types.Pointer:
		return Pointer(path, st, u)

	case *types.Basic:
		return Basic(path, st, u, named)

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

	sort.Strings(s)

	return strings.Trim(strings.Join(s, " | "), " ")
}

func makeIOutput(strukt *types.Struct, typename *types.TypeName) IObj {
	iobj := dump(typename.Name(), strukt, reflect.StructTag("")).(IStruct)

	for i := range iobj.items {
		if item, ok := iobj.items[i].(IBasic); ok {
			switch item.name {
			case "kind":
				item.options = typename.Name()
				iobj.items[i] = item
			case "apiVersion":
				item.options = typename.Pkg().Name()
				iobj.items[i] = item
			}
		}
	}

	return iobj
}

func pathAppend(path string, a string) string {
	if a == "" {
		return path
	} else if path == "" {
		return a
	} else {
		return path + "." + a
	}
}
