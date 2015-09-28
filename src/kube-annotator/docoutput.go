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
	"golang.org/x/tools/go/types"
	"reflect"
	"strings"
)

type DocOutput struct {
	indent string
}

func (o *DocOutput) Struct(basename *types.TypeName, st reflect.StructTag, u *types.Struct, options string, tn string) {
	name := strings.Split(st.Get("json"), ",")[0]

	if name != "" && name != basename.Name() {
		print(o.indent, name + ":", "", tn, st.Get("description"))
	}
	
	for i := 0; i < u.NumFields(); i++ {
		if name != "" && basename.Name() != name {
			o.indent = o.indent + "  "
		}

		dump(o, basename, u.Field(i).Type(), reflect.StructTag(u.Tag(i)))

		if name != "" && basename.Name() != name {
			o.indent = o.indent[:len(o.indent) - 2]
		}
		
		o.indent = strings.Replace(o.indent, "-", " ", -1)
	}
}

func (o *DocOutput) Map(basename *types.TypeName, st reflect.StructTag, u *types.Map, options string, tn string) {
	name := strings.Split(st.Get("json"), ",")[0]

	print(o.indent, name + ":", options, tn, st.Get("description"))
	print(o.indent + "  ", "[" + typefmt(u.Key()) + "]:", "", typefmt(u.Elem()), "")
}

func (o *DocOutput) Slice(basename *types.TypeName, st reflect.StructTag, u *types.Slice, options string) {
	name := strings.Split(st.Get("json"), ",")[0]

	und := u.Elem()
	if p, ok := und.(*types.Pointer); ok {
		und = p.Elem()
	}
	
	print(o.indent, name + ":", options, "[]" + typefmt(und), st.Get("description"))
	
	if _, ok := und.Underlying().(*types.Struct); ok {
		o.indent = o.indent + "- "
		dump(o, basename, und, reflect.StructTag("json:\"\""))
		o.indent = o.indent[:len(o.indent) - 2]
	} else {
		print(o.indent + "- ", "[" + typefmt(und) + "]", "", "", "")
	}
}

func (o *DocOutput) Pointer(basename *types.TypeName, st reflect.StructTag, u *types.Pointer) {
	dump(o, basename, u.Elem(), st)
}

func (o *DocOutput) Basic(basename *types.TypeName, st reflect.StructTag, u *types.Basic, options string, tn string) {
	name := strings.Split(st.Get("json"), ",")[0]

	if o.indent == "" {
		switch name {
		case "kind":
			options = basename.Name()
		case "apiVersion":
			options = basename.Pkg().Name()
		}
	}
	
	print(o.indent, name + ":", options, tn, st.Get("description"))
}

func print(indent, s1, options, s2, desc string) {
	if desc != "" {
		desc = " (" + desc + ")"
	}
	if options != "" {
		options = " " + options
	}

	fmt.Printf("%-69s # %s\n", indent + s1 + options, s2 + desc)
}
