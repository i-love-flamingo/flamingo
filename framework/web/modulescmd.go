package web

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"

	"flamingo.me/flamingo/v3/framework/config"
)

// ModulesCmd for debugging the router configuration
func ModulesCmd(area *config.Area) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "modules",
		Short: "Modules dump",
		Run: func(_ *cobra.Command, _ []string) {

			dumpModules(area)

		},
	}

	return cmd
}

func dumpModules(area *config.Area) {
	if area == nil {
		return
	}
	fmt.Println()
	fmt.Println("****************************************************************************")
	fmt.Println("Modules Names:               ")
	fmt.Println("****************************************************************************")

	for _, module := range area.Modules {
		moduleType := reflect.TypeOf(module)
		if moduleType.Kind() == reflect.Ptr {
			fmt.Print("*")
		}
		fmt.Printf("%s.%s\n", moduleType.Elem().PkgPath(), moduleType.Elem().Name())
	}
}
