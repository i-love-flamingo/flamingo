package commands

import (
	"github.com/spf13/cobra"
	"flamingo/framework/context"
	"flamingo/core/cmd/application"
	"fmt"
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
		fmt.Println("\nRouting Result:\n")
		routers := application.GetRouterForRootContext()
		for baseUrl,router := range routers {
			fmt.Printf("%s:\n   ",baseUrl)
			for routeName,route := range router.GetRoutes() {
				handler,_ := router.GetHandleForNamedRoute(routeName)
				fmt.Printf("    %s: %s: %B\n",routeName,route,handler)
			}
		}
	},
}
