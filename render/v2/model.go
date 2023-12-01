package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

var models = make(map[string]*Schema)

func registerModel(name string, schema *Schema) string {
	ref := "#/definitions/" + name
	models[name] = schema
	return ref
}

func renderSchema(obj spec.DefineStruct) (string, *Schema) {
	obj, ok := findType(obj.RawName) // go zero 解析路由时，结构体的 Member 会缺失，找到原始定义
	if !ok {
		panic("unknown type: " + obj.RawName)
	}
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
	} else {
		schema = renderPrimitiveProperty(field)
	}
	if schema == nil {
		return Property{}
	}
	schema.required = !isOptionalTag(tag)
	return Property{Name: tag.Name, Schema: schema}
}

func _renderPrimitiveProperty(typ spec.PrimitiveType) *Schema {
	typS, format := rawTypeFormat(typ.Name())
	if typS == "" {
		return nil
	}
	return &Schema{
		Type:   typS,
		Format: format,
	}
}

func renderPrimitiveProperty(field spec.Member) *Schema {
	schema := _renderPrimitiveProperty(field.Type.(spec.PrimitiveType))
	schema.Description = parseComment(field.Comment)
	return schema
}

func renderArrayProperty(array spec.ArrayType) *Schema {
	schema := &Schema{
		Type:        "array",
		Format:      "",
		Description: "",
	}
	memberTyp := array.Value
	if stru, ok := asDefineStruct(memberTyp); ok {
		_, items := renderSchema(stru)
		schema.Items = items
	} else if mArray, ok := asArrayType(memberTyp); ok {
		items := renderArrayProperty(mArray)
		schema.Items = items
	} else {
		items := _renderPrimitiveProperty(memberTyp.(spec.PrimitiveType))
		schema.Items = items
	}
	return schema
}
