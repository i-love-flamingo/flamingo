package command

import (
	"flamingo/framework/context"
	"fmt"

	"github.com/spf13/cobra"
)

var (
	contextName, baseUrl string
	Root                 *context.Context

	// DiCmd shows dependency injection information
	DiCmd = &cobra.Command{
		Use:   "di",
		Short: "Dependency Injection Debug output (for all or selected contexts)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nContainer for Routed Contexts:\n")
			for _, routeConfig := range Root.GetFlatContexts() {
				if contextName != "" && contextName != routeConfig.Name {
					continue
				}
				if baseUrl != "" && baseUrl != routeConfig.BaseURL {
					continue
				}
				fmt.Println()
				fmt.Println("********************************************")
				fmt.Println("Routed Context  - Baseurl:" + routeConfig.BaseURL + " Contextpath: [" + routeConfig.Name + "]")
				routeConfig.Injector.Debug()
				fmt.Println()
			}
		},
	}
)

func init() {
	DiCmd.Flags().StringVarP(&contextName, "context", "c", "", "Name of the context (relative context path) - set this if you like to see only this context. Otherwise it will show all.")
	DiCmd.Flags().StringVarP(&baseUrl, "baseurl", "", "", "Baseurl assigned to the context  - set this if you like to see only this context. Otherwise it will show all.")
}
