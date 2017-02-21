package template

import (
	"flamingo/core/flamingo"
	"net/http"
)

var TemplateFunctions = new(TplFuncRegistry)

func Register(basedir string) flamingo.RegisterFunc {
	return func(serviceContainer *flamingo.ServiceContainer) {
		serviceContainer.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir(basedir))))
		serviceContainer.Route("/static/{n:.*}", "_static")

		serviceContainer.Handle("_pugtpl_debug", new(DebugController))
		serviceContainer.Route("/_pugtpl/debug", "_pugtpl_debug")

		serviceContainer.Register(NewPugTemplateEngine(basedir))
		serviceContainer.Register(TemplateFunctions)
		serviceContainer.Register(new(AssetFunc), "template.func")
		serviceContainer.Register(new(DebugFunc), "template.func")
	}
}
