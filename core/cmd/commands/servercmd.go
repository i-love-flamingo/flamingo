package commands

import (
	"github.com/spf13/cobra"
	"net/http"
	"flamingo/core/cmd/application"
)


func init() {
	RootCommand.AddCommand(ServerCmd)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the main Web Server for the project",
	Run: func(cmd *cobra.Command, args []string) {
		http.ListenAndServe(":3210", application.GetFrontRouterForRootContext())
	},
}
