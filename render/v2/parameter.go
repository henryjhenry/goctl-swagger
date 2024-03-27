package v2

import (
	"github.com/henryjhenry/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

type member struct {
	m   spec.Member
	tag *spec.Tag
}

type members struct {
	pathMembers  []member // path parameters
	queryMembers []member // query parameters
	bodyMembers  []member // request body
}

// flatten flattens the struct and make it easier to render
func flatten(obj spec.DefineStruct) *members {
	var (
		pathMembers  []member // path parameters
		queryMembers []member // query parameters
		bodyMembers  []member // request body
	)
	for _, field := range obj.Members {
		if field.Name == "" {
			stru, _ := asDefineStruct(field.Type)
			subMembers := flatten(mustFindType(stru.RawName))
			pathMembers = append(pathMembers, (*subMembers).pathMembers...)
			queryMembers = append(queryMembers, (*subMembers).queryMembers...)
			bodyMembers = append(bodyMembers, (*subMembers).bodyMembers...)
			continue
		}
		tags := field.Tags()
		if len(tags) == 0 {
			continue
		}
		tag := lookupGozeroTag(tags)
		if tag == nil {
			continue
		}
		switch tag.Key {
		case types.PathTagKey: // path parameter only support primitive type
			if _, ok := field.Type.(spec.PrimitiveType); !ok {
				continue
			}
			pathMembers = append(pathMembers, member{m: field, tag: tag})
		case types.FormTagKey: // query parameter only support primitive type
			if _, ok := field.Type.(spec.PrimitiveType); !ok {
				continue
			}
			queryMembers = append(queryMembers, member{m: field, tag: tag})
		case types.JsonTagKey:
			bodyMembers = append(bodyMembers, member{m: field, tag: tag})
		}
	}

	return &members{pathMembers: pathMembers, queryMembers: queryMembers, bodyMembers: bodyMembers}
}

func renderParameters(obj spec.DefineStruct, method string) []*Parameter {
	members := flatten(obj)
	params := make([]*Parameter, 0, len((*members).pathMembers)+len((*members).queryMembers)+1)
	pathParams := renderPathParameters((*members).pathMembers)
	queryParams := renderQueryParameters((*members).queryMembers)
	params = append(append(params, pathParams...), queryParams...)
	body := renderRequestBody(obj.Name(), (*members).bodyMembers)
	if body != nil {
		params = append(params, body)
	}
	return params
}

func renderPathParameters(members []member) []*Parameter {
	if len(members) == 0 {
		return nil
	}
	parameters := make([]*Parameter, len(members))
	for i, m := range members {
		var param Parameter
		param.In = "path"
		param.Name = m.tag.Name
		param.Description = parseComment(m.m.Comment)
		param.Type, param.Format = rawTypeFormat(m.m.Type.Name())
		defaultVal, hasDefault := defaultTag(m.tag)
		if hasDefault {
			param.Required = false
			param.Default = defaultVal
		} else {
			param.Required = !isOptionalTag(m.tag)
		}
		parameters[i] = &param
	}
	return parameters
}

func renderQueryParameters(members []member) []*Parameter {
	if len(members) == 0 {
		return nil
	}
	parameters := make([]*Parameter, len(members))
	for i, m := range members {
		var param Parameter
		param.In = "query"
		param.Name = m.tag.Name
		param.Description = parseComment(m.m.Comment)
		param.Type, param.Format = rawTypeFormat(m.m.Type.Name())
		defaultVal, hasDefault := defaultTag(m.tag)
		if hasDefault {
			param.Required = false
			param.Default = defaultVal
		} else {
			param.Required = !isOptionalTag(m.tag)
		}
		parameters[i] = &param
	}
	return parameters
}

func renderRequestBody(name string, members []member) *Parameter {
	if len(members) == 0 {
		return nil
	}
	param := Parameter{
		Name: "body",
		In:   "body",
	}
	var (
		schema        = &Schema{Type: "object"}
		props         = make(Properties, 0, len(members))
		requiredProps []string
	)
	for _, m := range members {
		var prop Property
		if stru, ok := asDefineStruct(m.m.Type); ok {
			name, sc := renderSchema(stru)
			prop = Property{Name: m.tag.Name, Schema: &Schema{Ref: registerModel(name, sc)}}
		} else if array, ok := asArrayType(m.m.Type); ok {
			sc := renderArrayProperty(array)
			sc.Description = parseComment(m.m.Comment)
			prop = Property{Name: m.tag.Name, Schema: sc}
		} else {
			sc := renderPrimitiveProperty(m.m)
			if sc == nil {
				continue
			}
			prop = Property{Name: m.tag.Name, Schema: sc}
		}
		if !isOptionalTag(m.tag) {
			requiredProps = append(requiredProps, prop.Name)
		}
		props = append(props, prop)
	}
	schema.Properties = props
	schema.Required = requiredProps
	param.Schema = &Schema{Ref: registerModel(name, schema)}
	return &param
}
