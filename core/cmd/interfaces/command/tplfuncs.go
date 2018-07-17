package command

import (
	"fmt"

	"reflect"

	"flamingo.me/flamingo/framework/template"
	"flamingo.me/flamingo/framework/web"
	"github.com/spf13/cobra"
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
				bu, _ := routeConfig.Configuration.Get("prefixrouter.baseurl")
				baseurl, ok := bu.(string)
				if !ok {
					continue
				}
				if baseURL != "" && baseURL != baseurl {
					continue
				}
				fmt.Println()
				fmt.Println("********************************************")
				fmt.Println("Routed Context  - Baseurl:" + baseurl + " Contextpath: [" + routeConfig.Name + "]")
				tfr := routeConfig.Injector.GetInstance(template.FunctionRegistry{}).(*template.FunctionRegistry)
				fmt.Println("Functions")
				for _, f := range tfr.templateFunctions {
					fmt.Printf("%s: %s (from %s)\n", f.Name(), reflect.ValueOf(f.Func()).String(), reflect.ValueOf(f).Type().String())
				}
				fmt.Println()
				fmt.Println("Context Functions")
				for _, f := range tfr.contextTemplateFunctions {
					fmt.Printf("%s: %s (from %s)\n", f.Name(), reflect.ValueOf(f.Func(web.NewContext())).String(), reflect.ValueOf(f).Type().String())
				}
				fmt.Println()
			}
		},
	}
)
