package v2

import (
	"fmt"

	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func asDefineStruct(obj spec.Type) (spec.DefineStruct, bool) {
	stru, ok := obj.(spec.DefineStruct)
	if ok {
		return stru, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asDefineStruct(ptr.Type)
	}
	return spec.DefineStruct{}, false
}

func asArrayType(obj spec.Type) (spec.ArrayType, bool) {
	array, ok := obj.(spec.ArrayType)
	if ok {
		return array, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asArrayType(ptr.Type)
	}
	return spec.ArrayType{}, false
}

func asPrimitiveType(obj spec.Type) (spec.PrimitiveType, bool) {
	primitive, ok := obj.(spec.PrimitiveType)
	if ok {
		return primitive, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asPrimitiveType(ptr.Type)
	}
	return spec.PrimitiveType{}, false
}

func asMapType(obj spec.Type) (spec.MapType, bool) {
	mapType, ok := obj.(spec.MapType)
	if ok {
		return mapType, ok
	}
	if ptr, ok := obj.(spec.PointerType); ok {
		return asMapType(ptr.Type)
	}
	return spec.MapType{}, false
}

func panicUnsupportType(typ spec.Type) {
	panic(fmt.Sprintf("unsopport type %T", typ))
}
