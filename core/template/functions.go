package template

import (
	"flamingo/core/flamingo/web"
	"html/template"
	"log"
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
		TemplateFunctions        func() []Function        `inject:""`
		ContextTemplateFunctions func() []ContextFunction `inject:""`
		ContextAware             map[string]ContextAware
	}
)

// Populate Template Registry, mapping short method names to Functions
func (tfr *FunctionRegistry) Populate() template.FuncMap {
	tfr.ContextAware = make(map[string]ContextAware)
	funcMap := make(template.FuncMap)

	for _, tplFunc := range tfr.TemplateFunctions() {
		if tplFunc != nil {
			log.Println(tplFunc.Name())
			funcMap[tplFunc.Name()] = tplFunc.Func()
		}
	}
	for _, tplFunc := range tfr.ContextTemplateFunctions() {
		if tplFunc != nil {
			log.Println(tplFunc.Name())
			funcMap[tplFunc.Name()] = tplFunc.Func
			tfr.ContextAware[tplFunc.Name()] = tplFunc.Func
		}
	}

	return funcMap
}
