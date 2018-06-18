package command

import (
	"encoding/json"
	"fmt"

	"flamingo.me/flamingo/framework/config"

	"github.com/spf13/cobra"
)

var (
	// ConfigCmd shows config information
	ConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Config dump",
		Run: func(cmd *cobra.Command, args []string) {
			if contextName != "" {
				for _, c := range ConfigArea.Flat() {
					if c.Name == contextName {
						ConfigArea = c
						break
					}
				}
			}

			if len(args) > 0 {
				for _, c := range args {
					cfg, _ := ConfigArea.Config(c)
					x, _ := json.MarshalIndent(cfg, "", "  ")
					fmt.Println(c)
					fmt.Println(string(x))
					fmt.Println()
				}
			} else {
				dumpConfigArea(ConfigArea)
			}
		},
	}
)

func init() {
	ConfigCmd.Flags().StringVarP(&contextName, "context", "c", "", "Name of the context (relative context path) - set this if you like to see only this context. Otherwise it will show all.")
}

func dumpConfigArea(a *config.Area) {
	fmt.Println("**************************")
	fmt.Println(a.Name)
	fmt.Println("**************************")
	fmt.Println()
	x, _ := json.MarshalIndent(a.Configuration, "", "  ")
	fmt.Println(string(x))
	for _, routeConfig := range a.Childs {
		dumpConfigArea(routeConfig)
	}
}
