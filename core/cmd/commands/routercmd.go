package commands

import (
	"github.com/spf13/cobra"
	"flamingo/framework/context"
	"flamingo/core/cmd/application"
	"fmt"
	"reflect"
	"strings"
	"math"
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
			fmt.Println(routeConfig.BaseURL+" ["+routeConfig.Name+"]")
			for _, route := range routeConfig.Routes {
				fmt.Printf("  * %s > %s \n",route.Path,route.Controller)
				fmt.Printf("      Args: %v \n",route.Args)
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

		for baseUrl,router := range routers {
			fmt.Printf("%s:\n",baseUrl)

			for routeName,route := range router.GetRoutes() {
				handler,_ := router.GetHandleForNamedRoute(routeName)

				spaceAmount1 := int(math.Max(0, float64(20-len(routeName))))
				spaceAmount2 := int(math.Max(0, float64(30-len(route))))
				//Basti - possible to get struct again (real controller?) or should we add a "hint" param to the Handle method during registration?
				fmt.Printf("    %s:%s%s%s(%s)\n",routeName,strings.Repeat(" ",spaceAmount1),route,strings.Repeat(" ",spaceAmount2),reflect.TypeOf(handler))
			}
		}
	},
}
