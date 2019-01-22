package command

import (
	"encoding/json"
	"fmt"

	"flamingo.me/flamingo/v3/framework/config"
	"github.com/spf13/cobra"
)

func ConfigCmd(area *config.Area) *cobra.Command {
	var contextName string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config dump",
		Run: func(cmd *cobra.Command, args []string) {
			if contextName != "" {
				for _, c := range area.Flat() {
					if c.Name == contextName {
						area = c
						break
					}
				}
			}

			if len(args) > 0 {
				for _, c := range args {
					cfg, _ := area.Config(c)
					x, _ := json.MarshalIndent(cfg, "", "  ")
					fmt.Println(c + ":")
					fmt.Println(string(x))
					fmt.Println()
				}
			} else {
				dumpConfigArea(area)
			}
		},
	}

	cmd.Flags().StringVarP(
		&contextName,
		"context",
		"c",
		"",
		"Name of the context (relative context path) - set this if you like to see only this context. Otherwise it will show all.",
	)

	return cmd
}

func dumpConfigArea(a *config.Area) {
	fmt.Println()
	fmt.Println("**************************")
	fmt.Println("Area: ", a.Name)
	fmt.Println("**************************")
	x, _ := json.MarshalIndent(a.Configuration, "", "  ")
	fmt.Println(string(x))
	for _, routeConfig := range a.Childs {
		dumpConfigArea(routeConfig)
	}
}
