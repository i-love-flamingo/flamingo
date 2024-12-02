package flamingo

import (
	"bytes"

	"github.com/spf13/cobra"
)

func VersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			var buffer bytes.Buffer

			appInfo := GetAppInfo()
			PrintAppInfo(&buffer, appInfo)

			cmd.Println(buffer.String())
		},
	}
}
