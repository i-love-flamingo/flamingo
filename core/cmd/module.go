package cmd

import (
	"go.aoe.com/flamingo/core/cmd/interfaces/command"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"

	"github.com/spf13/cobra"
)

// Name for the default command
var Name = "flamingo"
var _ dingo.Module = &Module{}

// Module for core/cmd
type Module struct {
	Root *config.Area `inject:""`
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	var rootCmd = &cobra.Command{
		Use:   Name,
		Short: Name + " Console",
		Long:  "The " + Name + " command line interfaces.",
	}

	command.ConfigArea = m.Root

	rootCmd.AddCommand(
		command.VersionCmd,
		command.DiCmd,
		command.RoutingConfCmd,
		command.RouterCmd,
		command.DataControllerCmd,
		command.ConfigCmd,
		command.TplfuncsCmd,
	)

	injector.Bind((*cobra.Command)(nil)).AnnotatedWith("flamingo").ToInstance(rootCmd)
}
