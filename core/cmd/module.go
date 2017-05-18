package cmd

import (
	"flamingo/core/cmd/interfaces/command"
	"flamingo/framework/context"
	"flamingo/framework/dingo"

	"github.com/spf13/cobra"
)

var Name = "flamingo"

type Module struct {
	Root *context.Context `inject:""`
}

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
