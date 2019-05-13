package config

import (
	"flamingo.me/dingo"
	"github.com/spf13/pflag"
)

type (
	// Module defines a dingo module which automatically binds provided config.
	// Normaly this module is not included in your flamingo projects bootstrap.
	//
	// Its can be useful for testing dingo.Module that require certain configuration to be set before. E.g.:
	//
	// cfgModule := &config.Module{
	// 		Map: config.Map{
	// 			"redirects.useInRouter":       true,
	// 			"redirects.useInPrefixRouter": true,
	// 		},
	// 	}
	//
	// 	if err := dingo.TryModule(cfgModule, module); err != nil {
	// 		t.Error(err)
	// 	}
	Module struct {
		Map
	}

	// Flags handles the persistent flags provided by the config module
	Flags struct{}
)

var flagSet *pflag.FlagSet

func init() {
	flagSet = pflag.NewFlagSet("flamingo.config", pflag.ContinueOnError)
	flagSet.BoolVar(&debugLog, "flamingo-config-log", false, "enable flamingo config loader logging")
	flagSet.StringArrayVar(&additionalConfig, "flamingo-config", []string{}, "add multiple flamingo config additions")
}

// Configure the Module
func (m *Module) Configure(injector *dingo.Injector) {
	for k, v := range m.Flat() {
		if v == nil {
			continue
		}
		injector.Bind(v).AnnotatedWith("config:" + k).ToInstance(v)
	}
}

// Configure DI
func (f *Flags) Configure(injector *dingo.Injector) {
	injector.BindMulti((*pflag.FlagSet)(nil)).ToInstance(flagSet)
}
