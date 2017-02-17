package template

import (
	"encoding/json"
	"flamingo/core/flamingo"
	"flamingo/core/flamingo/web"
	"flamingo/core/packages/pug-template/pug-ast"
	"html/template"
	"strings"
)

type (
	TplFunc interface {
		Name() string
		Func() interface{}
	}

	ContextAware func(ctx web.Context) interface{}

	TplFuncRegistry struct {
		ServiceContainer *flamingo.ServiceContainer `inject:""`
		contextaware     map[string]ContextAware
	}

	AssetFunc struct {
		Router *flamingo.Router `inject:""`
	}
	DebugFunc struct{}
)

func (tfr *TplFuncRegistry) Populate() {
	tfr.contextaware = make(map[string]ContextAware)

	for _, tplfunc := range tfr.ServiceContainer.GetByTag("pug-template.func") {
		if tplfunc, ok := tplfunc.(TplFunc); ok {
			node.FuncMap[tplfunc.Name()] = tplfunc.Func()
			if f, ok := tplfunc.Func().(func(ctx web.Context) interface{}); ok {
				tfr.contextaware[tplfunc.Name()] = f
			}
		}
	}
}

func (_ AssetFunc) Name() string {
	return "asset"
}

func (af *AssetFunc) Func() interface{} {
	return func(a string) template.URL {
		if webpackserver {
			return template.URL("/assets/" + a)
		}

		url := af.Router.Url("_static", "n", "")
		aa := strings.Split(a, "/")
		aaa := aa[len(aa)-1]
		var result string
		if assetrewrites[aaa] != "" {
			result = url.String() + "/" + assetrewrites[aaa]
		} else {
			result = url.String() + "/" + a
		}
		//ctx.Push(result, nil)
		return template.URL(result)
	}
}

func (_ DebugFunc) Name() string {
	return "debug"
}

func (_ DebugFunc) Func() interface{} {
	return func(o interface{}) string {
		d, _ := json.MarshalIndent(o, "", "    ")
		return string(d)
	}
}
