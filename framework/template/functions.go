package template

import (
	"html/template"

	"flamingo.me/flamingo/framework/web"
)

type (
	// Function is a function which will be available in templates
	Function interface {
		Name() string
		Func() interface{}
	}

	// ContextFunction is a Function with late context binding
	ContextFunction interface {
		Name() string
		Func(web.Context) interface{}
	}

	// ContextAware is the used for late-bindings
	ContextAware func(ctx web.Context) interface{}

	// FunctionRegistry knows about the context-aware template functions
	FunctionRegistry struct {
		templateFunctions        []Function
		contextTemplateFunctions []ContextFunction
		ContextAware             map[string]ContextAware
	}
)

func (tfr *FunctionRegistry) Inject(templateFunctions []Function, contextTemplateFunctions []ContextFunction) {
	tfr.templateFunctions = templateFunctions
	tfr.contextTemplateFunctions = contextTemplateFunctions
}

// Populate Template Registry, mapping short method names to Functions
func (tfr *FunctionRegistry) Populate() template.FuncMap {
	tfr.ContextAware = make(map[string]ContextAware)
	funcMap := make(template.FuncMap)

	for _, tplFunc := range tfr.templateFunctions {
		if tplFunc != nil {
			funcMap[tplFunc.Name()] = tplFunc.Func()
		}
	}
	for _, tplFunc := range tfr.contextTemplateFunctions {
		if tplFunc != nil {
			funcMap[tplFunc.Name()] = tplFunc.Func
			tfr.ContextAware[tplFunc.Name()] = tplFunc.Func
		}
	}

	return funcMap
}
