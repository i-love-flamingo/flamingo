package application

import (
	"flamingo/core/prefix_router"
	"flamingo/framework/context"
	"flamingo/framework/router"
	"io"
	"log"
	"net/http"
	"os"
)

func GetFrontRouterForRootContext() *prefix_router.FrontRouter {
	return GetInitializedFrontRouter(context.RootContext.GetRoutingConfigs())
}

func GetRouterForRootContext() map[string]*router.Router {
	result := make(map[string]*router.Router)
	for _, routeConfig := range context.RootContext.GetRoutingConfigs() {
		result[routeConfig.BaseURL] = getRouterInContext(routeConfig)
	}
	return result
}

func GetInitializedFrontRouter(routingConfigs []*context.RoutingConfig) *prefix_router.FrontRouter {
	frontRouter := prefix_router.NewFrontRouter()
	defaultRouter := http.NewServeMux()
	frontRouter.Default(defaultRouter)
	addDefaultRoutes(defaultRouter)
	for _, routeConfig := range routingConfigs {
		routeConfig.Injector.Bind(new(log.Logger)).ToInstance(log.New(os.Stdout, "["+routeConfig.Name+"] ", 0))
		log.Println(routeConfig.Name, "at", routeConfig.BaseURL)
		frontRouter.Add(routeConfig.BaseURL, getRouterInContext(routeConfig))

	}
	return frontRouter
}

func getRouterInContext(routeConfig *context.RoutingConfig) *router.Router {
	return routeConfig.Injector.GetInstance(router.Router{}).(*router.Router).Init(routeConfig)
}

// TODO - Move this to a package?
func addDefaultRoutes(defaultRouter *http.ServeMux) {
	defaultRouter.HandleFunc("/assets/", func(rw http.ResponseWriter, req *http.Request) {
		if r, e := http.Get("http://localhost:1337" + req.RequestURI); e == nil {
			io.Copy(rw, r.Body)
		} else {
			rw.WriteHeader(404)
		}
	})

	defaultRouter.HandleFunc("/ping", func(rw http.ResponseWriter, req *http.Request) {
		rw.WriteHeader(200)
		rw.Write([]byte("pong"))
	})

	defaultRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Location", "/de/")
		rw.WriteHeader(301)
	})
}
