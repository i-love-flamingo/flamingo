package command

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"flamingo/framework/router"

	"flamingo/framework/web"

	"github.com/spf13/cobra"
)

type routesHelper struct {
	RouterRegistry *router.Registry `inject:""`
}

var (

	// RoutingConfCmd to show routing configuration information
	RoutingConfCmd = &cobra.Command{
		Use:   "routeconf",
		Short: "Print the routing configs from the contexts",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nContext with Routing Config:\n")
			for _, routeConfig := range ConfigArea.GetFlatContexts() {
				fmt.Println(routeConfig.BaseURL + " [" + routeConfig.Name + "]")
				for _, route := range routeConfig.Routes {
					fmt.Printf("  * %s > %s \n", route.Path, route.Controller)
				}
			}
		},
	}

	// RoutingConfCmd to show routing configuration information
	RouterCmd = &cobra.Command{
		Use:   "routes",
		Short: "Print the routes registered",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nRoutes:\n")
			RoutesHelper := ConfigArea.GetInitializedInjector().GetInstance(routesHelper{}).(*routesHelper)
			RoutesHelper.PrintRoutes()
		},
	}

	// RoutingConfCmd to show routing configuration information
	DataControllerCmd = &cobra.Command{
		Use:   "datacontroller",
		Short: "Print the datacontroller handlers registered. Datacontrollers can be called in Templates and also via Ajax",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Datacontroller:\n")
			RoutesHelper := ConfigArea.GetInitializedInjector().GetInstance(routesHelper{}).(*routesHelper)
			RoutesHelper.PrintDataHandlers()
		},
	}
)

// Print Registered Routes and Theire Handle
func (r *routesHelper) PrintRoutes() {
	routes := make(map[string]string)

	for _, routeHandler := range r.RouterRegistry.GetRoutes() {
		routes[routeHandler.GetPath()] = routeHandler.GetHandlerName()
	}

	fmt.Println("    Route-Name:            Route-Path                 (Registered Handler)")
	fmt.Println("--------------------------------------------------------------------------")

	for _, routePath := range getSortedMapKeys(routes) {
		_, controller := r.RouterRegistry.GetControllerForHandle(routes[routePath])
		printRoute(routes[routePath], routePath, controller)
	}
}

// Print Registered Routes and Theire Handle
func (r *routesHelper) PrintDataHandlers() {
	fmt.Println("    Handler-Name:         Type        (Registered Handler)")
	fmt.Println("----------------------------------------------------------")

	for k, v := range r.RouterRegistry.GetHandler() {
		if c, ok := v.(router.DataController); ok {
			fmt.Printf("    %s:\t\t> %s \t(%s)\n", k, "DataController", c)
		}
		if c, ok := v.(func(web.Context) interface{}); ok {
			fmt.Printf("    %s:\t\t> %s \t(%s)\n", k, "Function", c)
		}
	}
}

func getSortedMapKeys(theMap map[string]string) []string {
	var keys []string
	for k := range theMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printRoute(routeName string, routePath string, handler interface{}) {
	spaceAmount1 := int(math.Max(0, float64(25-len(routePath))))
	spaceAmount2 := int(math.Max(0, float64(30-len(routeName))))
	var handlerOutput string
	switch handler.(type) {
	case string:
		handlerOutput = handler.(string)
	default:
		handlerOutput = fmt.Sprintf("%T", handler)
	}
	fmt.Printf("    %s:%s%s%s(%s)\n", routePath, strings.Repeat(" ", spaceAmount1), routeName, strings.Repeat(" ", spaceAmount2), handlerOutput)
}
