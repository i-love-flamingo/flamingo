package commands

import (
	"github.com/spf13/cobra"
	"net/http"
	"flamingo/core/cmd/application"
	"fmt"
)


func init() {
	RootCommand.AddCommand(ServerCmd)
}

var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs the main Web Server for the project",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting HTTP Server at :3210 .....")
		e := http.ListenAndServe(":3210", application.GetFrontRouterForRootContext())
		if e != nil {
			fmt.Printf("Unexpected Error: %s", e)
		}
	},
}
