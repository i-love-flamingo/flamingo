package commands

import (
	"github.com/spf13/cobra"
	"flamingo/framework/context"
	"fmt"
)

func init() {
	RootCommand.AddCommand(DiCmd)
}


var DiCmd = &cobra.Command{
	Use:   "di",
	Short: "Dependency Injection",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\nContainer for Routed Contexts:\n")
		for _, routeConfig := range context.RootContext.GetRoutingConfigs() {
			fmt.Println(routeConfig.BaseURL+" ["+routeConfig.Name+"]")
			fmt.Println("********************************************")
			container := routeConfig.ServiceContainer
			container.Debug()
		}
	},
}
