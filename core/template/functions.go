package template

import (
	di "flamingo/core/flamingo/dependencyinjection"
	"flamingo/core/flamingo/web"
	"html/template"
)

type (
	// TemplateFunction is a function which will be available in templates
	TemplateFunction interface {
		Name() string
		Func() interface{}
	}

	// TemplateContextFunction is a TemplateFunction with late context binding
	TemplateContextFunction interface {
		Name() string
		Func(web.Context) interface{}
	}

	// ContextAware is the used for late-bindings
	ContextAware func(ctx web.Context) interface{}

	// TemplateFunctionRegistry knows about the context-aware template functions
	TemplateFunctionRegistry struct {
		ServiceContainer *di.Container `inject:""`
		Contextaware     map[string]ContextAware
	}
)

// Populate Template Registry, mapping short method names to Functions
func (tfr *TemplateFunctionRegistry) Populate() template.FuncMap {
	tfr.Contextaware = make(map[string]ContextAware)
	funcmap := make(template.FuncMap)

	for _, tplfunc := range tfr.ServiceContainer.GetTagged("template.func") {
		if tplfunc, ok := tplfunc.Value.(TemplateFunction); ok {
			funcmap[tplfunc.Name()] = tplfunc.Func()
		}
		if tplfunc, ok := tplfunc.Value.(TemplateContextFunction); ok {
			funcmap[tplfunc.Name()] = tplfunc.Func
			tfr.Contextaware[tplfunc.Name()] = tplfunc.Func
		}
	}

	return funcmap
}
