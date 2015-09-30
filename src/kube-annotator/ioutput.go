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
	"golang.org/x/tools/go/types"
	"reflect"
	"strings"
)

type IOutput struct {
}

type IBase struct {
	name string
	typ string
	description string
}

type IStruct struct {
	IBase
	items []interface{}
}

type IMap struct {
	IBase
	keytyp string
	valtyp string
}

type ISlice struct {
	IBase
	items []interface{}
	valtyp string
}

type IBasic struct {
	IBase
	options string
}

func (o *IOutput) Struct(st reflect.StructTag, u *types.Struct, tn string) interface{} {
	rv := IStruct{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = tn
	rv.description = st.Get("description")

	for i := 0; i < u.NumFields(); i++ {
		rv2 := dump(o, u.Field(i).Type(), reflect.StructTag(u.Tag(i)))
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

func (o *IOutput) Map(st reflect.StructTag, u *types.Map, tn string) interface{} {
	rv := IMap{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = tn
	rv.description = st.Get("description")
	rv.keytyp = typefmt(u.Key())
	rv.valtyp = typefmt(u.Elem())

	return rv
}

func (o *IOutput) Slice(st reflect.StructTag, u *types.Slice) interface{} {
	und := u.Elem()
	if p, ok := und.(*types.Pointer); ok {
		und = p.Elem()
	}

	rv := ISlice{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = "[]" + typefmt(und)
	rv.description = st.Get("description")
	c := dump(o, und, reflect.StructTag("json:\"\""))
	if st, ok := c.(IStruct); ok {
		rv.items = st.items
	} else {
		rv.valtyp = typefmt(und)
	}

	return rv
}

func (o *IOutput) Pointer(st reflect.StructTag, u *types.Pointer) interface{} {
	return dump(o, u.Elem(), st)
}

func (o *IOutput) Basic(st reflect.StructTag, u *types.Basic, options string, tn string) interface{} {
	rv := IBasic{}
	rv.name = strings.Split(st.Get("json"), ",")[0]
	rv.typ = tn
	rv.description = st.Get("description")
	rv.options = options

	return rv
}
