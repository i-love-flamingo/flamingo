package application

import (
	"os"
	"flamingo/framework/router"
	"net/http"
	"io"
	"flamingo/core/prefix_router"
	"flamingo/framework/context"
	"log"
)

func GetFrontRouterForRootContext() *prefix_router.FrontRouter  {
	return GetInitializedFrontRouter(context.RootContext.GetRoutingConfigs())
}

func GetRouterForRootContext() map[string]*router.Router  {
	result := make(map[string]*router.Router)
	for _, routeConfig := range context.RootContext.GetRoutingConfigs() {
		result[routeConfig.BaseURL] = router.CreateRouter(routeConfig)
	}
	return result
}

func GetInitializedFrontRouter(routingConfigs []*context.RoutingConfig) *prefix_router.FrontRouter {
	frontRouter := prefix_router.NewFrontRouter()
	defaultRouter := http.NewServeMux()
	frontRouter.Default(defaultRouter)
	addDefaultRoutes(defaultRouter)
	for _, routeConfig := range routingConfigs {
		// Register logger
		routeConfig.ServiceContainer.Register(log.New(os.Stdout, "["+routeConfig.Name+"] ", 0))
		//log.Println(routeConfig.Name, "at", routeConfig.BaseURL)
		frontRouter.Add(routeConfig.BaseURL, router.CreateRouter(routeConfig))
	}
	return frontRouter
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

