package web

import (
	"fmt"
	"reflect"

	"flamingo.me/dingo"
	"github.com/spf13/cobra"

	"flamingo.me/flamingo/v3/framework/config"
)

const showAllDuplicatesArg = "-a"

// ModulesCmd for debugging the router configuration
func ModulesCmd(area *config.Area) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "modules",
		Short: "Modules dump",
		Run: func(_ *cobra.Command, args []string) {

			argsMap := make(map[string]struct{})
			for _, arg := range args {
				argsMap[arg] = struct{}{}
			}

			dumpModules(area, argsMap)

		},
	}

	return cmd
}

func dumpModules(area *config.Area, argsMap map[string]struct{}) {
	if area == nil {
		return
	}

	_, printDuplicatedModules := argsMap[showAllDuplicatesArg]

	fmt.Println()
	fmt.Println("****************************************************************************")
	fmt.Println("Root Area Modules:               ")
	fmt.Println("****************************************************************************")

	registry := make(map[string]struct{})

	for _, module := range area.Modules {
		moduleName := getModuleName(module)
		registry[moduleName] = struct{}{}
		fmt.Print(moduleName)
	}

	for _, childArea := range area.Childs {
		fmt.Println()
		fmt.Println("****************************************************************************")
		fmt.Printf("Child Area: %s\n", childArea.Name)
		fmt.Println("****************************************************************************")

		for _, module := range childArea.Modules {
			moduleName := getModuleName(module)
			_, foundDuplicate := registry[moduleName]
			printModuleName := !foundDuplicate || printDuplicatedModules
			if printModuleName {
				fmt.Print(moduleName)
			}
		}
	}
	fmt.Println()
}

func getModuleName(module dingo.Module) string {
	moduleType := reflect.TypeOf(module)
	toRemember := ""
	if moduleType.Kind() == reflect.Ptr {
		toRemember = toRemember + "*"
	}
	return fmt.Sprintf("%s%s.%s\n", toRemember, moduleType.Elem().PkgPath(), moduleType.Elem().Name())
}
