package pug_template

import (
	"flamingo/core/flamingo/service_container"
	"flamingo/core/packages/pug_template/pugast"
	"flamingo/core/packages/pug_template/template_functions"
	"flamingo/core/template"
	"net/http"
)

func Register(basedir string, debug bool) service_container.RegisterFunc {
	return func(serviceContainer *service_container.ServiceContainer) {
		serviceContainer.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(basedir))))
		serviceContainer.Route("/static/{n:.*}", "_static")

		serviceContainer.Handle("_pugtpl_debug", new(DebugController))
		serviceContainer.Route("/_pugtpl/debug", "_pugtpl_debug")

		serviceContainer.Register(pugast.NewPugTemplateEngine(basedir, debug))
		serviceContainer.Register(new(template.TemplateFunctionRegistry))
		serviceContainer.Register(new(template_functions.AssetFunc), "template.func")
		serviceContainer.Register(new(template_functions.DebugFunc), "template.func")
	}
}
