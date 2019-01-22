package command

import (
	"fmt"

	"flamingo.me/flamingo/v3/framework/config"

	"github.com/spf13/cobra"
)

var (
	contextName, baseURL string

	// ConfigArea is the Root config area
	ConfigArea *config.Area

	// DiCmd shows dependency injection information
	DiCmd = &cobra.Command{
		Use:   "di",
		Short: "Dependency Injection Debug output (for all or selected contexts)",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nContainer for Routed Contexts:")
			fmt.Println()
			for _, routeConfig := range ConfigArea.GetFlatContexts() {
				if contextName != "" && contextName != routeConfig.Name {
					continue
				}
				bu, _ := routeConfig.Configuration.Get("prefixrouter.baseurl")
				baseurl, _ := bu.(string)
				if baseURL != "" && baseURL != baseurl {
					continue
				}
				fmt.Println()
				fmt.Println("********************************************")
				fmt.Println("Routed Context  - Baseurl:" + baseurl + " Contextpath: [" + routeConfig.Name + "]")
				routeConfig.Injector.Debug()
				fmt.Println()
			}
		},
	}
)

func init() {
	DiCmd.Flags().StringVarP(&contextName, "context", "c", "", "Name of the context (relative context path) - set this if you like to see only this context. Otherwise it will show all.")
	DiCmd.Flags().StringVarP(&baseURL, "baseurl", "", "", "Baseurl assigned to the context  - set this if you like to see only this context. Otherwise it will show all.")
}
