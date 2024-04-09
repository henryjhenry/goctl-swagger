package v2

import (
	"github.com/henryjhenry/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var models = make(map[string]*Schema)

func registerModel(name string, schema *Schema) string {
	ref := "#/definitions/" + name
	if _, ok := models[ref]; ok {
		return ref // do not register twice
	}
	models[name] = schema
	return ref
}

func renderSchema(obj spec.DefineStruct) (string, *Schema) {
	obj = mustFindType(obj.RawName) // go zero 解析路由时，结构体的 Member 会缺失，找到原始定义
	schema := &Schema{
		Type: "object",
	}
	properties := make(Properties, 0, len(obj.Members))
	var requiredProps []string
	for _, field := range obj.Members {
		if field.Name == "" { // 匿名字段一定是结构体
			stru, ok := asDefineStruct(field.Type)
			if !ok {
				continue
			}
			_, s := renderSchema(stru)
			for _, v := range s.Properties {
				properties = append(properties, v)
				if v.Schema.required {
					requiredProps = append(requiredProps, v.Name)
				}
			}
			continue
		}
		prop := renderProperty(field)
		if prop.Schema == nil {
			continue
		}
		properties = append(properties, prop)
		if prop.Schema.required {
			requiredProps = append(requiredProps, prop.Name)
		}
	}
	schema.Properties = properties
	if len(requiredProps) > 0 {
		schema.Required = requiredProps
	}
	return obj.Name(), schema
}

func renderProperty(field spec.Member) Property {
	tags := field.Tags()
	if len(tags) == 0 {
		return Property{}
	}
	tag := lookupGozeroTag(tags)
	if tag == nil || tag.Key != types.JsonTagKey {
		return Property{}
	}
	var schema *Schema
	typ := field.Type
	if stru, ok := asDefineStruct(typ); ok {
		_, schema = renderSchema(stru)
	} else if array, ok := asArrayType(typ); ok {
		schema = renderArrayProperty(array)
		schema.Description = parseComment(field.Comment) // reset description
	} else if primitive, ok := asPrimitiveType(typ); ok {
		schema = renderPrimitiveProperty(primitive)
		if schema == nil {
			panicUnsupportType(typ)
		}
		schema.Description = parseComment(field.Comment)
	} else if mapType, ok := asMapType(typ); ok {
		schema = renderMapTypeProperty(mapType)
		schema.Description = parseComment(field.Comment)
	}
	if schema == nil {
		panicUnsupportType(typ)
	}
	schema.required = !isOptionalTag(tag)
	return Property{Name: tag.Name, Schema: schema}
}

func renderMapTypeProperty(mapType spec.MapType) *Schema {
	if mapType.Key != "string" {
		panic("just support string type key on map")
	}
	valueType, ok := asPrimitiveType(mapType.Value)
	if !ok {
		panic("just support primitive type value on map")
	}
	valtyp, _ := primitiveTypeFormat(valueType.Name())
	return &Schema{
		Type: "object",
		AdditionalProperties: &Schema{
			Type: valtyp,
		},
	}
}

func renderPrimitiveProperty(typ spec.PrimitiveType) *Schema {
	typS, format := primitiveTypeFormat(typ.Name())
	if typS == "" {
		return nil
	}
	return &Schema{
		Type:   typS,
		Format: format,
	}
}

func renderPrimitivePropertyByMember(field spec.Member) *Schema {
	schema := renderPrimitiveProperty(field.Type.(spec.PrimitiveType))
	if schema != nil {
		schema.Description = parseComment(field.Comment)
	}
	return schema
}

func renderArrayProperty(array spec.ArrayType) *Schema {
	var (
		schema = &Schema{
			Type:        "array",
			Format:      "",
			Description: "",
		}
		memberTyp = array.Value
		items     *Schema
	)

	if stru, ok := asDefineStruct(memberTyp); ok {
		_, items = renderSchema(stru)
	} else if mArray, ok := asArrayType(memberTyp); ok {
		items = renderArrayProperty(mArray)
	} else if primitive, ok := asPrimitiveType(memberTyp); ok {
		items = renderPrimitiveProperty(primitive)
	}
	if items == nil {
		panicUnsupportType(memberTyp)
	}
	schema.Items = items
	return schema
}
