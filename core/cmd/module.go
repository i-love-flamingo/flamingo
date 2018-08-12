package cmd

import (
	"os"
	"path/filepath"

	"flamingo.me/flamingo/core/cmd/interfaces/command"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"github.com/spf13/cobra"
)

// Module for DI
type Module struct{}

var dingoTrace *bool

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	//command.VersionCmd,
	//command.DiCmd,
	//command.RoutingConfCmd,
	//command.RouterCmd,
	//command.DataControllerCmd,
	//command.TplfuncsCmd,

	injector.Bind(new(cobra.Command)).AnnotatedWith("flamingo").ToProvider(
		func(commands []*cobra.Command, config *struct {
			Name string `inject:"config:cmd.name"`
		}) *cobra.Command {
			rootCmd := &cobra.Command{
				Use:   config.Name,
				Short: "Flamingo " + config.Name,
			}

			rootCmd.AddCommand(commands...)

			return rootCmd
		},
	)

	injector.BindMulti(new(cobra.Command)).ToProvider(command.ConfigCmd)
}

// DefaultConfig specifies the command name
func (m *Module) DefaultConfig() config.Map {
	return config.Map{
		"cmd.name": filepath.Base(os.Args[0]),
	}
}

// Run the root command
func Run(injector *dingo.Injector) error {
	return injector.GetAnnotatedInstance(new(cobra.Command), "flamingo").(*cobra.Command).Execute()
}
