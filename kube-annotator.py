#!/usr/bin/python

import collections
import json
import re
import sys


def loadapi(fn):
    with open(fn) as f:
        return json.load(f, object_pairs_hook=collections.OrderedDict)


def resolvetype(p):
    if "type" in p:
        typ = p["type"]

        if typ == "array":
            typ = "[]" + resolvetype(p["items"])
            
    elif "$ref" in p:
        typ = p["$ref"]

    else:
        raise Exception()

    return typ


def printproperties(api, m, indent=""):
    for p in api["models"][m]["properties"]:
        typ = resolvetype(api["models"][m]["properties"][p])
        if typ == "object":
            typ = "map[string]string"

        key = p + ":"
        if indent == "":
            if p == "kind":
                key += " " + m.split(".", 1)[1]
            elif p == "apiVersion":
                key += " " + m.split(".", 1)[0]
                
        print "%-40s# %-40s%s" % (indent + key, typ, re.sub(r"\n+", " ", api["models"][m]["properties"][p].get("description", "")))

        indent = indent.replace("-", " ")

        rtyp = typ.lstrip("[]")
        if rtyp in api["models"]:
            if typ == rtyp:
                newindent = indent + "  "
            else:
                newindent = indent + "- "
            printproperties(api, rtyp, newindent)
        elif typ == "map[string]string":
            print "%-40s# %s" % (indent + "  [string]:", "string")
        elif rtyp != typ:
            print "%-40s# %s" % (indent + "- [" + rtyp + "]", "")
       

def printmodels(api):
    for m in sorted(api["models"]):
        if m[0] != "v":
            continue

        if "kind" not in api["models"][m]["properties"]:
            continue

        if m.endswith("List"):
            continue
        
        print m.split(".", 1)[1]
        print "=" * len(m.split(".", 1)[1])
        print

        if "description" in api["models"][m]:
            print api["models"][m]["description"]
            print

        printproperties(api, m)
        print
        print


def main(fn):
    printmodels(loadapi(fn))


if __name__ == "__main__":
    main(*sys.argv[1:])
