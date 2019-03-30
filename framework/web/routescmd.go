package web

import (
	"fmt"
	"math"
	"sort"
	"strings"

	"flamingo.me/flamingo/v3/framework/config"
	"github.com/spf13/cobra"
)

// RoutesCmd for debugging the router configuration
func RoutesCmd(router *Router, area *config.Area) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "routes",
		Short: "Routes dump",
		Run: func(cmd *cobra.Command, args []string) {

			dumpRoutes(router, area)

		},
	}

	return cmd
}

// HandlerCmd for debugging the router/handler configuration
func HandlerCmd(router *Router, area *config.Area) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "handler",
		Short: "Dump the Handlers and its registered methods",
		Run: func(cmd *cobra.Command, args []string) {

			dumpHandler(router, area)

		},
	}

	return cmd
}

func dumpRoutes(router *Router, area *config.Area) {
	// router.Init(area)
	fmt.Println()
	fmt.Println("***************************************************************************")
	fmt.Println(" Route                						| Handler-Name:               ")
	fmt.Println("****************************************************************************")
	for _, routeHandler := range router.routerRegistry.routes {
		routePath := routeHandler.path.path + "(" + strings.Join(routeHandler.path.params, ";") + ")"
		spaceAmount1 := int(math.Max(0, float64(60-len(routePath))))
		fmt.Printf("    %s%s| %s\n", routePath, strings.Repeat(" ", spaceAmount1), routeHandler.handler)
	}
}

func dumpHandler(router *Router, area *config.Area) {
	// router.Init(area)
	fmt.Println()
	fmt.Println("***************************************************************************")
	fmt.Println(" Handle-name                	 | registered actions               ")
	fmt.Println("****************************************************************************")

	handlerNamesSorted := getSortedMapKeys(router.routerRegistry.handler)
	for _, handlerKey := range handlerNamesSorted {
		handler := router.routerRegistry.handler[handlerKey]
		var actions []string
		if handler.data != nil {
			actions = append(actions, "DATA")
		}
		if handler.data != nil {
			actions = append(actions, "ANY")
		}
		for method := range handler.method {
			actions = append(actions, method)
		}
		spaceAmount1 := int(math.Max(0, float64(30-len(handlerKey))))

		fmt.Printf(" %s %s | %s   \n", handlerKey, strings.Repeat(" ", spaceAmount1), strings.Join(actions, " ; "))
	}
}

func getSortedMapKeys(theMap map[string]handlerAction) []string {
	var keys []string
	for k := range theMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
