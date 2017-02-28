package template

import (
	"flamingo/core/flamingo/service_container"
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
		ServiceContainer *service_container.ServiceContainer `inject:""`
		Contextaware     map[string]ContextAware
	}
)

// Populate Template Registry, mapping short method names to Functions
func (tfr *TemplateFunctionRegistry) Populate() template.FuncMap {
	tfr.Contextaware = make(map[string]ContextAware)
	funcmap := make(template.FuncMap)

	for _, tplfunc := range tfr.ServiceContainer.GetByTag("template.func") {
		if tplfunc, ok := tplfunc.(TemplateFunction); ok {
			funcmap[tplfunc.Name()] = tplfunc.Func()
		}
		if tplfunc, ok := tplfunc.(TemplateContextFunction); ok {
			funcmap[tplfunc.Name()] = tplfunc.Func
			tfr.Contextaware[tplfunc.Name()] = tplfunc.Func
		}
	}

	return funcmap
}
