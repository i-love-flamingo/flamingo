package command

import (
	"flamingo/framework"
	"fmt"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flamingo",
	Long:  `All software has versions. This is Flamingo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Flamingo Console v%s -- HEAD\n", framework.VERSION)
	},
}
