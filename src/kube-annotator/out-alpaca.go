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
	"os"

	"golang.org/x/tools/go/types"
)

type JsonObject map[string]interface{}

func write(filename string, jsobj JsonObject) {
	s, _ := json.MarshalIndent(jsobj, "", "  ")
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(s)
}

func alpacaf(strukt *types.Struct, typename *types.TypeName, docpkg *doc.Package) string {
	schema := makeSchema(makeIOutput(strukt, typename))
	delete(schema["properties"].(JsonObject), "status")
	write("js/schema.json", schema)

	options := makeOptions(makeIOutput(strukt, typename))
	delete(options["fields"].(JsonObject), "status")
	write("js/options.json", options)

	return ""
}

func makeSchema(iobj IObj) JsonObject {
	switch iobj := iobj.(type) {
	case IStruct:
		jsobj := make(JsonObject)
		if iobj.items != nil {
			jsobj["type"] = "object"
			properties := make(JsonObject)
			for _, item := range iobj.items {
				childjsobj := makeSchema(item)
				if childjsobj != nil {
					properties[item.Name()] = childjsobj
				}
			}
			jsobj["properties"] = properties
		} else {
			jsobj["type"] = "string"
		}

		return jsobj

	case IMap:
		jsobj := make(JsonObject)
		jsobj["type"] = "array"

		items := make(JsonObject)
		items["type"] = "object"

		properties := make(JsonObject)

		_key := make(JsonObject)
		_key["title"] = "_key"
		_key["type"] = "string"
		properties["_key"] = _key

		val := make(JsonObject)
		val["title"] = "val"
		val["type"] = "string"
		properties["val"] = val

		items["properties"] = properties

		jsobj["items"] = items

		return jsobj

	case ISlice:
		// TODO: rbd.monitors...
		jsobj := make(JsonObject)
		jsobj["type"] = "array"

		items := make(JsonObject)
		items["type"] = "object"

		properties := make(JsonObject)
		for _, item := range iobj.items {
			childjsobj := makeSchema(item)
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

		return jsobj
	}

	return nil
}

func makeOptions(iobj IObj) JsonObject {
	switch iobj := iobj.(type) {
	case IStruct:
		jsobj := make(JsonObject)
		jsobj["label"] = iobj.name
		//jsobj["collapsed"] = true
		if iobj.items != nil {
			fields := make(JsonObject)
			for i, item := range iobj.items {
				childjsobj := makeOptions(item)
				if childjsobj != nil {
					childjsobj["order"] = i
					fields[item.Name()] = childjsobj
				}
			}
			jsobj["fields"] = fields
		} else {
			jsobj["type"] = "text"
		}

		return jsobj

	case IMap:
		jsobj := make(JsonObject)
		jsobj["label"] = iobj.name
		jsobj["type"] = "map"

		items := make(JsonObject)
		items["type"] = "object"

		fields := make(JsonObject)

		_key := make(JsonObject)
		_key["label"] = "key"
		_key["type"] = "text"
		fields["_key"] = _key

		val := make(JsonObject)
		val["label"] = "val"
		val["type"] = "text"
		fields["val"] = val

		items["fields"] = fields

		jsobj["items"] = items

		return jsobj

	case ISlice:
		jsobj := make(JsonObject)
		jsobj["label"] = iobj.name
		//jsobj["collapsed"] = true

		items := make(JsonObject)

		fields := make(JsonObject)
		for i, item := range iobj.items {
			childjsobj := makeOptions(item)
			if childjsobj != nil {
				childjsobj["order"] = i
				fields[item.Name()] = childjsobj
			}
		}
		items["fields"] = fields
		jsobj["items"] = items

		return jsobj

	case IBasic:
		jsobj := make(JsonObject)
		jsobj["label"] = iobj.name
		jsobj["type"] = "text"

		return jsobj
	}

	return nil
}
