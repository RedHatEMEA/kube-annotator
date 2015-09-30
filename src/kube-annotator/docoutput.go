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
	"strings"
)

func do2(o interface{}) string {
	rv := ""
	switch u := o.(type) {
	case IStruct:
		for _, i := range u.items {
			rv += do3("", i)
		}
	}

	return rv
}

func do3(indent string, o interface{}) string {
	rv := ""
	switch u := o.(type) {
	case IStruct:
		rv = sprint(indent, u.name + ":", "", u.typ, u.description)
		for _, i := range u.items {
			rv += do3(indent + "  ", i)
		}

	case IMap:
		rv = sprint(indent, u.name + ":", "", u.typ, u.description)
		rv += sprint(indent + "  ", "[" + u.keytyp + "]:", "", u.valtyp, "")

	case ISlice:
		rv = sprint(indent, u.name + ":", "", u.typ, u.description)
		rv2 := ""
		if len(u.items) > 0 {
			for i := range u.items {
				rv2 += do3(indent + "  ", u.items[i])
			}
			rv += strings.Replace(rv2, indent + "  ", indent + "- ", 1)
		} else {
			rv += sprint(indent + "- ", "[" + u.valtyp + "]", "", "", "")
		}

	case IBasic:
		rv = sprint(indent, u.name + ":", u.options, u.typ, u.description)
	}

	return rv
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

	return fmt.Sprintf("%-69s #%s\n", indent + s1 + options, s2 + desc)
}
