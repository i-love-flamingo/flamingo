// Package ops is a test module with instance routes, a command, a template function, a config schema and a legacy config alias.
package ops

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Module of the ops routes and commands.
	Module struct{}

	// SingletonInstance is bound in dingo.Singleton.
	// The field gives the type a size, because Go gives one address to all zero-size allocations.
	SingletonInstance struct{ _ int }

	// ChildInstance is bound in dingo.ChildSingleton.
	// The field gives the type a size, because Go gives one address to all zero-size allocations.
	ChildInstance struct{ _ int }

	routes struct {
		responder *web.Responder
		singleton singletonProvider
		child     childProvider
	}

	// upperFunc is the template function "upper".
	upperFunc struct{}

	// dingo makes a provider function for a type whose name ends with "Provider".
	singletonProvider func() *SingletonInstance
	childProvider     func() *ChildInstance
)

// Configure DI.
func (*Module) Configure(injector *dingo.Injector) {
	injector.Bind(SingletonInstance{}).In(dingo.Singleton)
	injector.Bind(ChildInstance{}).In(dingo.ChildSingleton)

	web.BindRoutes(injector, new(routes))

	injector.BindMulti(new(cobra.Command)).ToProvider(greetCommand)

	flamingo.BindTemplateFunc(injector, "upper", new(upperFunc))
}

// Func gives the function of the template function "upper".
func (*upperFunc) Func(context.Context) interface{} {
	return strings.ToUpper
}

// CueConfig gives the greeting of the ops-greet command.
func (*Module) CueConfig() string {
	return `ops: { greeting: string | *"hello" }`
}

// FlamingoLegacyConfigAlias maps the old key of the greeting to the current one.
func (*Module) FlamingoLegacyConfigAlias() map[string]string {
	return map[string]string{"opsGreeting": "ops.greeting"}
}

func greetCommand(cfg *struct {
	Greeting string `inject:"config:ops.greeting"`
}) *cobra.Command {
	return &cobra.Command{
		Use:   "ops-greet",
		Short: "print the configured greeting",
		Run: func(cmd *cobra.Command, _ []string) {
			cmd.Println(cfg.Greeting)
		},
	}
}

// Inject dependencies.
// The instances come through providers, so each request asks its scope again.
func (r *routes) Inject(responder *web.Responder, singleton singletonProvider, child childProvider) {
	r.responder = responder
	r.singleton = singleton
	r.child = child
}

// Routes of the module.
func (r *routes) Routes(registry *web.RouterRegistry) {
	registry.MustRoute("/ops/instance/singleton", "ops.instance.singleton")
	registry.HandleGet("ops.instance.singleton", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]string{"id": fmt.Sprintf("%p", r.singleton())})
	})

	registry.MustRoute("/ops/instance/child", "ops.instance.child")
	registry.HandleGet("ops.instance.child", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]string{"id": fmt.Sprintf("%p", r.child())})
	})
}
