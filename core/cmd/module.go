package cmd

import (
	"flamingo/core/cmd/interfaces/command"
	"flamingo/framework/config"
	"flamingo/framework/dingo"

	"github.com/spf13/cobra"
)

// Name for the default command
var Name = "flamingo"

// Module for core/cmd
type Module struct {
	Root *config.Area `inject:""`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	var rootCmd = &cobra.Command{
		Use:   Name,
		Short: "Flamingo Console",
		Long:  `The flamingo command line interface.`,
	}

	command.Root = m.Root

	rootCmd.AddCommand(
		command.VersionCmd,
		command.DiCmd,
		command.RoutingConfCmd,
		//command.RouterCmd,
	)

	injector.Bind((*cobra.Command)(nil)).AnnotatedWith("flamingo").ToInstance(rootCmd)
}
