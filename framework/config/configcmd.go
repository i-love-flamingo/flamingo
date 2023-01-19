package config

import (
	"encoding/json"
	"fmt"

	"cuelang.org/go/cue/format"
	"github.com/spf13/cobra"
)

// Cmd command: The Area for which the config is to be printed need to be passed. This will be done by Dingo if a Provider is used for example.
func Cmd(area *Area) *cobra.Command {
	var contextName string

	cmd := &cobra.Command{
		Use:   "config",
		Short: "Config dump",
		Run: func(cmd *cobra.Command, args []string) {
			if contextName != "" {
				flatArea, _ := area.Flat()
				for _, c := range flatArea {
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

func dumpConfigArea(a *Area) {
	fmt.Println()
	fmt.Println("**************************")
	fmt.Println("Area: ", a.Name)
	fmt.Println("**************************")
	if false { // cuedump {
		// build a cue runtime to verify the config
		// cueRuntime := new(cue.Runtime)
		//ci, err := cueRuntime.Build(a.cueBuildInstance)
		//if err != nil {
		//	panic(err)
		//}

		for _, f := range a.cueBuildInstance.Files {
			d, _ := format.Node(f, format.Simplify())
			fmt.Println("//", f.Filename)
			fmt.Println(string(d))
			fmt.Println("")
		}

		// d, _ := format.Node(ci.Value().Syntax(), format.Simplify())
		//fmt.Println(string(d))
	} else {
		x, _ := json.MarshalIndent(a.Configuration, "", "  ")
		fmt.Println(string(x))
	}
	for _, routeConfig := range a.Childs {
		dumpConfigArea(routeConfig)
	}
}
