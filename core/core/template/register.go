package template

import (
	"flamingo/core/core/app"
	"net/http"
)

var TemplateFunctions = new(TplFuncRegistry)

func Register(serviceContainer *app.ServiceContainer) {
	serviceContainer.Handle("_static", http.StripPrefix("/static/", http.FileServer(http.Dir("frontend/dist"))))
	serviceContainer.Route("/static/{n:.*}", "_static")

	serviceContainer.Register(TemplateFunctions)
	serviceContainer.Register(new(AssetFunc), "template.func")
	serviceContainer.Register(new(DebugFunc), "template.func")
}
