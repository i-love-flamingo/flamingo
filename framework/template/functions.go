package template

import (
	"context"

	"flamingo.me/flamingo/v3/framework/dingo"
)

type (
	Func interface {
		Func() interface{}
	}

	CtxFunc interface {
		Func(context.Context) interface{}
	}

	FuncProvider    func() map[string]Func
	CtxFuncProvider func() map[string]CtxFunc
)

func BindFunc(injector *dingo.Injector, name string, fnc Func) {
	injector.BindMap(new(Func), name).To(fnc)
}

func BindCtxFunc(injector *dingo.Injector, name string, fnc CtxFunc) {
	injector.BindMap(new(CtxFunc), name).To(fnc)
}
