// Package ops is a test module with instance and event routes, a command, a healthcheck status,
// a template function, a config schema and a legacy config alias.
package ops

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/core/healthcheck/domain/healthcheck"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// startupDelays makes StartupEvent slower than ServerStartEvent, so a server that does not wait for
// StartupEvent answers before it is handled.
const startupDelays = 2

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
		started   *started
		singleton singletonProvider
		child     childProvider
	}

	// started tells which start events the app has sent. It is a singleton, so the routes of an area and the
	// event subscriber of the root area share it.
	started struct {
		startup     atomic.Bool
		serverStart atomic.Bool
	}

	// eventLog logs the lifecycle events of the app and records the start events.
	eventLog struct {
		logger  flamingo.Logger
		started *started
		delay   time.Duration // wait before a start event is recorded
	}

	// status is the healthcheck status "ops". It fails when the config says so.
	status struct {
		fail bool
	}

	// twice is a result that writes the header two times.
	twice struct{}

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

	injector.Bind(new(started)).In(dingo.Singleton)

	flamingo.BindEventSubscriber(injector).To(new(eventLog))

	injector.BindMap(new(healthcheck.Status), "ops").To(new(status))
}

// Inject dependencies.
func (e *eventLog) Inject(logger flamingo.Logger, state *started, cfg *struct {
	DelayMS int `inject:"config:ops.events.delayMS"`
}) {
	e.logger = logger.WithField(flamingo.LogKeyCategory, "contract")
	e.started = state
	e.delay = time.Duration(cfg.DelayMS) * time.Millisecond
}

// Notify logs the type of each lifecycle event. The start events also set their flag, before the log line.
func (e *eventLog) Notify(_ context.Context, event flamingo.Event) {
	switch event.(type) {
	case *flamingo.StartupEvent:
		time.Sleep(startupDelays * e.delay)
		e.started.startup.Store(true)
	case *flamingo.ServerStartEvent:
		time.Sleep(e.delay)
		e.started.serverStart.Store(true)
	case *flamingo.ServerShutdownEvent, *flamingo.ShutdownEvent:
	default:
		return
	}

	e.logger.Info(fmt.Sprintf("contract event %T", event))
}

// Inject dependencies.
func (s *status) Inject(cfg *struct {
	Fail bool `inject:"config:ops.health.fail"`
}) {
	s.fail = cfg.Fail
}

// Status is down when "ops.health.fail" is true.
func (s *status) Status() (bool, string) {
	if s.fail {
		return false, "ops down"
	}

	return true, "ops up"
}

// Apply writes the status twice, so the http server logs "superfluous response.WriteHeader call".
func (twice) Apply(_ context.Context, w http.ResponseWriter) error {
	w.WriteHeader(http.StatusOK)
	w.WriteHeader(http.StatusOK)

	return nil
}

// Func gives the function of the template function "upper".
func (*upperFunc) Func(context.Context) interface{} {
	return strings.ToUpper
}

// CueConfig is the schema of the ops config.
func (*Module) CueConfig() string {
	return `ops: { greeting: string | *"hello", health: fail: bool | *false, events: delayMS: number | *0 }`
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
func (r *routes) Inject(responder *web.Responder, state *started, singleton singletonProvider, child childProvider) {
	r.responder = responder
	r.started = state
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

	registry.MustRoute("/ops/twice", "ops.twice")
	registry.HandleGet("ops.twice", func(_ context.Context, _ *web.Request) web.Result {
		return twice{}
	})

	registry.MustRoute("/ops/started", "ops.started")
	registry.HandleGet("ops.started", func(_ context.Context, _ *web.Request) web.Result {
		return r.responder.Data(map[string]bool{
			"startup":     r.started.startup.Load(),
			"serverStart": r.started.serverStart.Load(),
		})
	})
}
