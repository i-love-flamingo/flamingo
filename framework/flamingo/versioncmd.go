package flamingo

import (
	"fmt"

	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(AppVersion())
		},
	}
}
