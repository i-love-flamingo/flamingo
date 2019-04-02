package command

import (
	"fmt"

	"go.aoe.com/flamingo/framework"

	"github.com/spf13/cobra"
)

// VersionCmd shows the Flamingo version
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Flamingo",
	Long:  `All software has versions. This is Flamingo's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Flamingo Console v%s -- HEAD\n", framework.VERSION)
	},
}
