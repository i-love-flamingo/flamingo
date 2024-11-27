package flamingo

import (
	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(AppVersion())
		},
	}
}
