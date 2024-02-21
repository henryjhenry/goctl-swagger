package v2

import (
	"net/http"
	"strings"

	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
)

func renderPaths(svc spec.Service, opt option) Paths {
	var (
		paths   Paths
		tags    = svc.Name
		pathMap = make(map[string]*Path)
	)
	for _, grp := range svc.Groups {
		if value := grp.GetAnnotation("group"); len(value) > 0 {
			tags = value
		}
		if len(opt.TagPrefix) > 0 {
			tags = opt.TagPrefix + tags
		}
		for _, route := range grp.Routes {
			uri := grp.GetAnnotation("prefix") + route.Path
			if uri[0] != '/' {
				uri = "/" + uri
			}
			var (
				path  *Path
				isNew bool
			)
			if obj, ok := pathMap[uri]; ok {
				path = obj
			} else {
				isNew = true
				path = &Path{Path: uri}
			}
			op := Operation{
				Summary: strings.Trim(route.AtDoc.Text, `"`),
				Tags:    []string{tags},
			}
			if obj, ok := route.RequestType.(spec.DefineStruct); ok {
				op.Parameters = renderParameters(obj, strings.ToUpper(route.Method))
			}

			// Just support json response.
			if obj, ok := route.ResponseType.(spec.DefineStruct); ok {
				// root schema

				op.Responses = map[string]*Response{
					"200": {
						Description: "OK",
						Schema:      renderRootSchema(opt, obj),
					},
				}

			} else {
				op.Responses = map[string]*Response{
					"200": {
						Description: "OK",
					},
				}
			}
			switch strings.ToUpper(route.Method) {
			case http.MethodGet:
				path.Get = &op
			case http.MethodPost:
				path.Post = &op
			case http.MethodPut:
				path.Put = &op
			case http.MethodDelete:
				path.Delete = &op
			case http.MethodPatch:
				path.Patch = &op
			case http.MethodHead:
				path.Head = &op
			case http.MethodOptions:
				path.Options = &op
			}
			if isNew {
				paths = append(paths, path)
				pathMap[uri] = path
			}
		}
	}
	paths.Sort()
	return paths
}

func renderRootSchema(opt option, obj spec.DefineStruct) (root *Schema) {
	_, resp := renderSchema(obj)
	ref := registerModel(obj.Name(), resp)
	if opt.outsideSchema != nil {
		root = renderOutsideSchema(*opt.outsideSchema)
		root.Properties = append(root.Properties, Property{"data", &Schema{
			Type: "object",
			Ref:  ref,
		}})
		root.Description = strings.Join(obj.Docs, ",")
		return
	}
	root = &Schema{
		Description: strings.Join(obj.Docs, ","),
		Ref:         ref,
	}
	return
}

func renderOutsideSchema(outSide spec.DefineStruct) *Schema {
	schema := &Schema{
		Type: "object",
	}
	properties := make(Properties, 0, len(outSide.Members))
	var requiredProps []string
	for _, field := range outSide.Members {
		if field.Name == "" { // 匿名字段一定是结构体
			continue
		}
		if _, ok := field.Type.(spec.PrimitiveType); !ok {
			continue
		}
		tags := field.Tags()
		if len(tags) == 0 {
			continue
		}
		tag := lookupGozeroTag(tags)
		if tag == nil || tag.Key != types.JsonTagKey {
			continue
		}
		prop := Property{Name: tag.Name, Schema: renderPrimitiveProperty(field)}
		properties = append(properties, prop)
		if prop.Schema.required {
			requiredProps = append(requiredProps, prop.Name)
		}
	}
	schema.Properties = properties
	if len(requiredProps) > 0 {
		schema.Required = requiredProps
	}
	return schema
}
