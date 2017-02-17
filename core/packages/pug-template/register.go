package template

import (
	"flamingo/core/flamingo"
	"net/http"
)

var TemplateFunctions = new(TplFuncRegistry)

func Register(serviceContainer *flamingo.ServiceContainer) {
	serviceContainer.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/dist"))))
	serviceContainer.Route("/static/{n:.*}", "_static")

	serviceContainer.Register(TemplateFunctions)
	serviceContainer.Register(new(AssetFunc), "pug-template.func")
	serviceContainer.Register(new(DebugFunc), "pug-template.func")
}
