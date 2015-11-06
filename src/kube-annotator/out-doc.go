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
	"go/doc"
	"strings"

	"golang.org/x/tools/go/types"
)

func docf(strukt *types.Struct, typename *types.TypeName, docpkg *doc.Package) string {
	s := typename.Name() + "\n"
	s += strings.Repeat("=", len(typename.Name())) + "\n\n"

	for _, t := range docpkg.Types {
		if t.Name == typename.Name() {
			s += t.Doc + "\n"
			break
		}
	}

	iobj := makeIOutput(strukt, typename)
	for _, item := range iobj.(IStruct).items {
		s += makeDoc("", item)
	}

	return s + "\n"
}

func makeDoc(indent string, iobj IObj) string {
	s := ""
	switch iobj := iobj.(type) {
	case IStruct:
		s = sprint(indent, iobj.name+":", "", iobj.typ, iobj.description)
		for _, item := range iobj.items {
			s += makeDoc(indent+"  ", item)
		}

	case IMap:
		s = sprint(indent, iobj.name+":", "", iobj.typ, iobj.description)
		s += sprint(indent+"  ", "["+iobj.keytyp+"]:", "", iobj.valtyp, "")

	case ISlice:
		s = sprint(indent, iobj.name+":", "", iobj.typ, iobj.description)
		if len(iobj.items) > 0 {
			s2 := ""
			for _, item := range iobj.items {
				s2 += makeDoc(indent+"  ", item)
			}
			s += strings.Replace(s2, indent+"  ", indent+"- ", 1)
		} else if iobj.valtyp != "" {
			s += sprint(indent+"- ", "["+iobj.valtyp+"]", "", "", "")
		}

	case IBasic:
		s = sprint(indent, iobj.name+":", iobj.options, iobj.typ, iobj.description)
	}

	return s
}

func sprint(indent, s1, options, s2, desc string) string {
	if desc != "" {
		desc = " (" + desc + ")"
	}
	if options != "" {
		options = " " + options
	}
	if s2 != "" {
		s2 = " " + s2
	}

	return fmt.Sprintf("%-69s #%s\n", indent+s1+options, s2+desc)
}
