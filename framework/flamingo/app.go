package flamingo

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/spf13/cobra"
	"github.com/zemirco/memorystore"
	"go.aoe.com/flamingo/framework/config"
	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
)

type appmodule struct {
	Cmd    *cobra.Command `inject:"flamingo"`
	Root   *config.Area   `inject:""`
	Router *router.Router `inject:""`
}

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	sessionStore := memorystore.NewMemoryStore([]byte("flamingosecret"))
	sessionStore.MaxLength(1024 * 1024)
	injector.Bind((*sessions.Store)(nil)).ToInstance(sessionStore)

	a.Cmd.AddCommand(&cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			a.Router.Init(a.Root)
			http.ListenAndServe(":3322", a.Router)
		},
	})
}

func (a *appmodule) OverrideConfig(config.Map) config.Map {
	return config.Map{
		"flamingo.template.err404": "404",
		"flamingo.template.err503": "503",
	}
}

// App is a simple app-runner for flamingo
func App(root *config.Area, configdir string) {
	app := new(appmodule)
	root.Modules = append(root.Modules, app)
	if configdir == "" {
		configdir = "config"
	}
	config.Load(root, configdir)

	if err := app.Cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
