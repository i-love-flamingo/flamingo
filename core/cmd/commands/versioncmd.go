package commands

import (
	"github.com/spf13/cobra"
	"fmt"
)


func init() {
	RootCommand.AddCommand(VersionCmd)
}

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flamingo",
	Long:  `All software has versions. This is Flamingo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Flamingo Console v0.1 -- HEAD")
	},
}
