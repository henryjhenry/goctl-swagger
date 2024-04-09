package render

import (
	"testing"

	"github.com/henryjhenry/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

func TestRender20(t *testing.T) {
	apiPath := "../testdata/main.api"
	apiSpec, err := parser.Parse(apiPath)
	if err != nil {
		t.Fatal(err)
	}
	plg := &plugin.Plugin{
		Api: apiSpec,
		Dir: ".",
	}
	outsideSchema := "../testdata/api/outside_schema.api"
	opt := types.Option{
		Target:        "../swagger.json",
		Version:       "2.0",
		RenderType:    "json",
		TagPrefix:     "",
		OutsideSchema: outsideSchema,
		Host:          "http://127.0.0.1:8888",
	}
	if err := Render(plg, opt); err != nil {
		t.Fatal(err)
	}
}
