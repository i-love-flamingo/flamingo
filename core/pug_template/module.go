package pug_template

import (
	"flamingo/core/pug_template/pugast"
	"flamingo/core/pug_template/template_functions"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/framework/template"
	template_functions2 "flamingo/framework/template_functions"
	"flamingo/framework/web"
	"net/http"
)

type (
	// Module for framework/pug_template
	Module struct {
		RouterRegistry *router.Registry `inject:""`
		Basedir        string           `inject:"config:pug_template.basedir"`
	}

	// TemplateFunctionInterceptor to use fixtype
	TemplateFunctionInterceptor struct {
		template.ContextFunction
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(m.Basedir))))
	m.RouterRegistry.Route("/static/*n", "_static")

	m.RouterRegistry.Route("/_pugtpl/debug", "pugtpl.debug")
	m.RouterRegistry.Handle("pugtpl.debug", new(DebugController))

	// We bind the Template Engine to the ChildSingleton level (in case there is different config handling
	// We use the provider to make sure both are always the same injected type
	injector.Bind(pugast.PugTemplateEngine{}).In(dingo.ChildSingleton)
	injector.Bind((*template.Engine)(nil)).
		In(dingo.ChildSingleton).
		ToProvider(func(t *pugast.PugTemplateEngine, i *dingo.Injector) template.Engine { return (template.Engine)(t) })

	injector.BindMulti((*template.ContextFunction)(nil)).To(template_functions.AssetFunc{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.MathLib{})
	injector.BindMulti((*template.Function)(nil)).To(template_functions.DebugFunc{})

	injector.BindInterceptor((*template.ContextFunction)(nil), TemplateFunctionInterceptor{})
}

// Func interceptor
// we want to intercept the GetFunc() to make sure we convert the result via pugast.Fixtype
// This allows to cut the dependency from framework to pug_template module
func (t *TemplateFunctionInterceptor) Func(ctx web.Context) interface{} {
	if getfunc, ok := t.ContextFunction.(*template_functions2.GetFunc); ok {
		oGetFunc := getfunc.Func(ctx).(func(string, ...map[interface{}]interface{}) interface{})
		return func(what string, params ...map[interface{}]interface{}) interface{} {
			return pugast.Fixtype(oGetFunc(what, params...))
		}
	}
	return t.ContextFunction.Func(ctx)
}
