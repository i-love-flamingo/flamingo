package template

import (
	"context"
	"html/template"

	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/web"
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

	// Function is a function which will be available in templates
	// deprecated: use Func
	Function interface {
		Name() string
		Func() interface{}
	}

	// ContextFunction is a Function with late context binding
	// deprecated: use CtxFunc
	ContextFunction interface {
		Name() string
		Func(web.Context) interface{}
	}

	// ContextAware is the used for late-bindings
	// deprecated: use CtxFunc
	ContextAware func(ctx web.Context) interface{}

	// FunctionRegistry knows about the context-aware template functions
	// deprecated: use CtxFuncProvider()
	FunctionRegistry struct {
		templateFunctions        []Function
		contextTemplateFunctions []ContextFunction
		ContextAware             map[string]ContextAware
	}
)

func BindFunc(injector *dingo.Injector, name string, fnc Func) {
	injector.BindMap(new(Func), name).To(fnc)
}

func BindCtxFunc(injector *dingo.Injector, name string, fnc CtxFunc) {
	injector.BindMap(new(CtxFunc), name).To(fnc)
}

func (tfr *FunctionRegistry) Inject(templateFunctions []Function, contextTemplateFunctions []ContextFunction) {
	tfr.templateFunctions = templateFunctions
	tfr.contextTemplateFunctions = contextTemplateFunctions
}

// Populate Template Registry, mapping short method names to Functions
// deprecated: use CtxFuncProvider()
func (tfr *FunctionRegistry) Populate() template.FuncMap {
	tfr.ContextAware = make(map[string]ContextAware)
	funcMap := make(template.FuncMap)

	for _, tplFunc := range tfr.templateFunctions {
		if tplFunc != nil {
			funcMap[tplFunc.Name()] = tplFunc.Func()
		}
	}
	for _, tplFunc := range tfr.contextTemplateFunctions {
		tplFunc := tplFunc
		if tplFunc != nil {
			funcMap[tplFunc.Name()] = func(ctx context.Context) interface{} {
				return tplFunc.Func(web.ToContext(ctx))
			}
			tfr.ContextAware[tplFunc.Name()] = tplFunc.Func
		}
	}

	return funcMap
}
