package commands

import (
	"flamingo/core/cmd/application"
	"flamingo/framework/context"
	"fmt"
	"math"
	"strings"

	"sort"

	"github.com/spf13/cobra"
)

func init() {
	RootCommand.AddCommand(RoutingConfCmd)
	RootCommand.AddCommand(RouterCmd)
}

var RoutingConfCmd = &cobra.Command{
	Use:   "routeconf",
	Short: "Print the routing configs from the contexts",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nContext with Routing Config:\n")
		for _, routeConfig := range context.RootContext.GetRoutingConfigs() {
			fmt.Println(routeConfig.BaseURL + " [" + routeConfig.Name + "]")
			for _, route := range routeConfig.Routes {
				fmt.Printf("  * %s > %s \n", route.Path, route.Controller)
				fmt.Printf("      Args: %v \n", route.Args)
			}
		}
	},
}

var RouterCmd = &cobra.Command{
	Use:   "routes",
	Short: "Print the routes after evaluating the config. With the real resulting Controller",
	Run: func(cmd *cobra.Command, args []string) {
		routers := application.GetRouterForRootContext()

		fmt.Println("\nRouting Result:\n")
		fmt.Println("    Route-Name:            Route-Path                 (Registered Handler)")
		fmt.Println("--------------------------------------------------------------------------")

		for baseUrl, router := range routers {
			fmt.Printf("\n**********\nContext: \"%s\":\n", baseUrl)
			fmt.Println("  Hardroutes:")
			for _, route := range router.GetHardRoutes() {
				printRoute("--", route.Path, route.Controller, route.Args)
			}

			fmt.Println("  Registered Routes:")
			routes := router.RouterRegistry.GetRoutes()
			for _, routeName := range getSortedMapKeys(routes) {
				route := routes[routeName]
				handler, _ := router.RouterRegistry.GetHandleForNamedRoute(routeName)
				printRoute(routeName, route, handler, nil)
			}
		}
		fmt.Println()
	},
}

func getSortedMapKeys(theMap map[string]string) []string {
	var keys []string
	for k := range theMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func printRoute(routeName string, routePath string, handler interface{}, args interface{}) {
	spaceAmount1 := int(math.Max(0, float64(20-len(routeName))))
	spaceAmount2 := int(math.Max(0, float64(30-len(routePath))))
	var handlerOutput string
	switch handler.(type) {
	case string:
		handlerOutput = handler.(string)
	default:
		handlerOutput = fmt.Sprintf("%T", handler)
	}
	fmt.Printf("    %s:%s%s%s(%s [%s])\n", routeName, strings.Repeat(" ", spaceAmount1), routePath, strings.Repeat(" ", spaceAmount2), handlerOutput, args)
}
