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
	"encoding/json"
	"go/doc"
	"golang.org/x/tools/go/types"
)

type JsonObject map[string]interface{}

func jsonf(strukt *types.Struct, typename *types.TypeName, docpkg *doc.Package) string {
	js := makeJsonObject(makeIOutput(strukt, typename))

	obj := make(JsonObject)
	obj["schema"] = js["properties"]
	delete(obj["schema"].(JsonObject), "status")

	s, _ := json.MarshalIndent(obj, "", "  ")
	return string(s)
}

func makeJsonObject(iobj IObj) JsonObject {
	switch iobj := iobj.(type) {
	case IStruct:
		jsobj := make(JsonObject)
		jsobj["title"] = iobj.name
		if iobj.items != nil {
			jsobj["type"] = "object"
			properties := make(JsonObject)
			for _, item := range iobj.items {
				childjsobj := makeJsonObject(item)
				if childjsobj != nil {
					properties[item.Name()] = childjsobj
				}
			}
			jsobj["properties"] = properties
		} else {
			jsobj["type"] = "string"
		}

		return jsobj

	case ISlice:
		jsobj := make(JsonObject)
		jsobj["type"] = "array"

		items := make(JsonObject)
		items["title"] = iobj.name
		items["type"] = "object"

		properties := make(JsonObject)
		for _, item := range iobj.items {
			childjsobj := makeJsonObject(item)
			if childjsobj != nil {
				properties[item.Name()] = childjsobj
			}
		}
		items["properties"] = properties

		jsobj["items"] = items

		return jsobj

	case IBasic:
		jsobj := make(JsonObject)
		jsobj["type"] = "string"
		jsobj["title"] = iobj.name

		return jsobj
	}

	return nil
}
