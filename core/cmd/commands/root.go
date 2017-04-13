package commands

import (
	"github.com/spf13/cobra"
)

// encapsulates the cobra Command
type FlamingoRootCommand struct {
	cobraRootCmd *cobra.Command
}

// simple command struct that can be used by other packages to register own commands
// TODO
type FlamingoCommand struct {
	name string
	Short string
	Run func(args []string)
}

var RootCommand = &FlamingoRootCommand{
	&cobra.Command{
		Short: "Flamingo Console",
		Long:  `The flamingo command line interface.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

		},
	},
}


func (RootCommand *FlamingoRootCommand) AddFlamingoCommand(Command *FlamingoCommand) {
	newCobraCommand := &cobra.Command{

	}
	RootCommand.cobraRootCmd.AddCommand(newCobraCommand)
}

func (RootCommand *FlamingoRootCommand) AddCommand(command *cobra.Command) {
	RootCommand.cobraRootCmd.AddCommand(command)
}

func (RootCommand *FlamingoRootCommand) Execute() error {
	return RootCommand.cobraRootCmd.Execute()
}

func (RootCommand *FlamingoRootCommand) SetUse(use string) {
	RootCommand.cobraRootCmd.Use = use
}
