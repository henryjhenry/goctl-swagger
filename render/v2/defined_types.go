package v2

import (
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var registeredTypes = make(map[string]spec.DefineStruct)

func register(typ spec.DefineStruct) {
	registeredTypes[typ.RawName] = typ
}

func registerTypes(types []spec.Type) {
	for i := range types {
		if stru, ok := asDefineStruct(types[i]); ok {
			register(stru)
		}
	}
}

func findType(name string) (spec.DefineStruct, bool) {
	obj, ok := registeredTypes[name]
	return obj, ok
}

func mustFindType(name string) spec.DefineStruct {
	obj, ok := findType(name)
	if !ok {
		panic("unknown type: " + name)
	}
	return obj
}
