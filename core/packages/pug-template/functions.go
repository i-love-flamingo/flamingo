package template

import (
	"flamingo/core/flamingo"
	"flamingo/core/flamingo/web"
	"flamingo/core/packages/pug-template/pugast"
)

type (
	TplFunc interface {
		Name() string
		Func() interface{}
	}

	TplContextFunc interface {
		Name() string
		Func(web.Context) interface{}
	}

	ContextAware func(ctx web.Context) interface{}

	TplFuncRegistry struct {
		ServiceContainer *flamingo.ServiceContainer `inject:""`
		contextaware     map[string]ContextAware
	}
)

func (tfr *TplFuncRegistry) Populate() {
	tfr.contextaware = make(map[string]ContextAware)

	for _, tplfunc := range tfr.ServiceContainer.GetByTag("template.func") {
		if tplfunc, ok := tplfunc.(TplFunc); ok {
			pugast.FuncMap[tplfunc.Name()] = tplfunc.Func()
		}
		if tplfunc, ok := tplfunc.(TplContextFunc); ok {
			pugast.FuncMap[tplfunc.Name()] = tplfunc.Func
			tfr.contextaware[tplfunc.Name()] = tplfunc.Func
		}
	}
}
