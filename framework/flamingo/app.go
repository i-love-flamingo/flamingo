package flamingo

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"flamingo.me/flamingo/core/cmd"
	"flamingo.me/flamingo/framework/config"
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
	"github.com/spf13/cobra"
)

type (
	appmodule struct {
		Root   *config.Area   `inject:""`
		Router *router.Router `inject:""`
	}
	// AppShutdownEvent is dispatched on app shutdown
	AppShutdownEvent struct {
		AppModule dingo.Module
	}
)

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(cobra.Command)).ToInstance(&cobra.Command{
		Use: "serve",
		Run: func(cmd *cobra.Command, args []string) {
			a.handleShutdown()
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

func (a *appmodule) handleShutdown() {
	signals := make(chan os.Signal)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func(m *appmodule) {
		<-signals
		m.Router.EventRouterProvider().Dispatch(context.Background(), &AppShutdownEvent{AppModule: m})
	}(a)
}

// App is a simple app-runner for flamingo
func App(root *config.Area, configdir string) {
	app := new(appmodule)
	root.Modules = append(root.Modules, app)
	if configdir == "" {
		configdir = "config"
	}
	config.Load(root, configdir)

	if err := cmd.Run(root.Injector); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
