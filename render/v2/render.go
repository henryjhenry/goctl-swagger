package v2

import (
	"path/filepath"
	"strings"

	"github.com/henryjhenry/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

type Renderer struct {
}

type option struct {
	types.Option
	outsideSchema *spec.DefineStruct
}

func (r *Renderer) Render(plg *plugin.Plugin, opt types.Option) (types.Swagger, error) {
	registerTypes(plg.Api.Types)
	var contact *Contact
	if len(plg.Api.Info.Email) > 0 || len(plg.Api.Info.Author) > 0 {
		contact = &Contact{
			Name:  plg.Api.Info.Author,
			Email: plg.Api.Info.Email,
		}
	}
	info := Information{
		Title:       strings.Trim(plg.Api.Info.Properties["title"], `"`),
		Description: plg.Api.Info.Properties["desc"],
		Version:     plg.Api.Info.Properties["version"],
		Contact:     contact,
	}
	_opt := option{
		Option: opt,
	}
	var oscPath string // outside schema path
	if filepath.IsAbs(opt.OutsideSchema) {
		oscPath = opt.OutsideSchema
	} else {
		oscPath = filepath.Join(plg.Dir, opt.OutsideSchema)
	}
	// if err != nil {
	// 	return nil, err
	// }
	if opt.OutsideSchema != "" {
		stru, err := readOutsideSchema(oscPath)
		if err != nil {
			return nil, err
		}
		_opt.outsideSchema = stru
	}
	paths := renderPaths(plg.Api.Service, _opt)
	swagger := &Swagger{
		Swagger:  "2.0",
		Info:     info,
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Paths:    paths,
		Schemes:  opt.Schemes,
	}
	swagger.Definitions = models
	return swagger, nil
}

func readOutsideSchema(path string) (*spec.DefineStruct, error) {
	apiSpec, err := parser.Parse(path)
	if err != nil {
		return nil, err
	}
	if len(apiSpec.Types) == 0 {
		return nil, nil
	}

	typ := apiSpec.Types[0]
	if stru, ok := typ.(spec.DefineStruct); ok {
		return &stru, nil
	}
	return nil, nil
}
