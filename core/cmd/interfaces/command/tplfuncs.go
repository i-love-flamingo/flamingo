package command

import (
	"fmt"

	"github.com/spf13/cobra"
	"go.aoe.com/flamingo/framework/template"
	"reflect"
	"go.aoe.com/flamingo/framework/web"
)

var (
	TplfuncsCmd = &cobra.Command{
		Use:   "tplfuncs",
		Short: "Debug Template Functions",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("\nContainer for Routed Contexts:")
			fmt.Println()
			for _, routeConfig := range ConfigArea.GetFlatContexts() {
				if contextName != "" && contextName != routeConfig.Name {
					continue
				}
				if baseURL != "" && baseURL != routeConfig.BaseURL {
					continue
				}
				fmt.Println()
				fmt.Println("********************************************")
				fmt.Println("Routed Context  - Baseurl:" + routeConfig.BaseURL + " Contextpath: [" + routeConfig.Name + "]")
				tfr := routeConfig.Injector.GetInstance(template.FunctionRegistry{}).(*template.FunctionRegistry)
				fmt.Println("Functions")
				for _, f := range tfr.TemplateFunctions {
					fmt.Printf("%s: %s (from %s)\n", f.Name(), reflect.ValueOf(f.Func()).String(), reflect.ValueOf(f).Type().String())
				}
				fmt.Println()
				fmt.Println("Context Functions")
				for _, f := range tfr.ContextTemplateFunctions {
					fmt.Printf("%s: %s (from %s)\n", f.Name(), reflect.ValueOf(f.Func(web.NewContext())).String(), reflect.ValueOf(f).Type().String())
				}
				fmt.Println()
			}
		},
	}
)
