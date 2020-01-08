package flamingo

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/core/runtime"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"
)

// fmtErrorf shim to allow go 1.12 backwards compatibility
var fmtErrorf = fmt.Errorf

type (
	// Application contains a main flamingo application
	Application struct {
		configDir       string
		childAreas      []*config.Area
		area            *config.Area
		args            []string
		defaultContext  string
		eagerSingletons bool
		flagset         *flag.FlagSet
	}

	// ApplicationOption configures an Application
	ApplicationOption func(config *Application)
)

// ConfigDir configuration ApplicationOption
func ConfigDir(configdir string) ApplicationOption {
	return func(config *Application) {
		config.configDir = configdir
	}
}

// ChildAreas allows to define additional config areas for roots
func ChildAreas(areas ...*config.Area) ApplicationOption {
	return func(config *Application) {
		config.childAreas = areas
	}
}

// DefaultContext for flamingo to start with
func DefaultContext(name string) ApplicationOption {
	return func(config *Application) {
		config.defaultContext = name
	}
}

// SetEagerSingletons controls if eager singletons will be created
func SetEagerSingletons(enabled bool) ApplicationOption {
	return func(config *Application) {
		config.eagerSingletons = enabled
	}
}

// WithArgs sets the initial arguments different than os.Args[1:]
func WithArgs(args ...string) ApplicationOption {
	return func(config *Application) {
		config.args = args
	}
}

type eventRouterProvider func() flamingo.EventRouter

type arrayFlags []string

func (i *arrayFlags) String() string {
	return strings.Join(*i, ", ")
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

// NewApplication loads a new application for running the Flamingo application with the given modules, loaded configs etc
func NewApplication(modules []dingo.Module, options ...ApplicationOption) (*Application, error) {
	app := &Application{
		configDir:      "config",
		args:           os.Args[1:],
		defaultContext: "root",
	}

	for _, option := range options {
		option(app)
	}

	app.flagset = flag.NewFlagSet("flamingo", flag.ContinueOnError)
	dingoTraceCircular := app.flagset.Bool("dingo-trace-circular", false, "enable dingo circular tracing")
	flamingoConfigLog := app.flagset.Bool("flamingo-config-log", false, "enable flamingo config logging")
	flamingoConfigCueDebug := app.flagset.String("flamingo-config-cue-debug", "", "query the flamingo cue config loader (use . for root)")
	flamingoContext := app.flagset.String("flamingo-context", app.defaultContext, "set flamingo execution context")
	var flamingoConfig arrayFlags
	app.flagset.Var(&flamingoConfig, "flamingo-config", "add additional flamingo yaml config")
	dingoInspect := app.flagset.Bool("dingo-inspect", false, "inspect dingo")

	if err := app.flagset.Parse(app.args); err != nil && err != flag.ErrHelp {
		return nil, fmtErrorf("app: parsing arguments: %w", err)
	}

	if dingoTraceCircular != nil && *dingoTraceCircular {
		dingo.EnableCircularTracing()
	}

	modules = append([]dingo.Module{
		new(framework.InitModule),
		new(zap.Module),
		new(runtime.Module),
		new(cmd.Module),
	}, modules...)
	modules = append(modules, new(servemodule))

	root := config.NewArea("root", modules, app.childAreas...)

	configLoadOptions := []config.LoadOption{
		config.AdditionalConfig(flamingoConfig),
		config.DebugLog(*flamingoConfigLog),
		config.LegacyMapping(true, false),
	}

	if *flamingoConfigCueDebug != "" {
		printCue := func(b []byte, err error) {
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(string(b))
			os.Exit(-1)
		}
		if *flamingoConfigCueDebug == "." {
			configLoadOptions = append(configLoadOptions, config.CueDebug(nil, printCue))
		} else {
			configLoadOptions = append(configLoadOptions, config.CueDebug(strings.Split(*flamingoConfigCueDebug, "."), printCue))
		}
	}

	if err := config.Load(root, app.configDir, configLoadOptions...); err != nil {
		return nil, fmtErrorf("app: config load: %w", err)
	}

	areas, err := root.Flat()
	if err != nil {
		return nil, fmtErrorf("app: flat areas: %w", err)
	}

	var ok bool
	app.area, ok = areas[*flamingoContext]
	if !ok {
		return nil, fmtErrorf("app: context %q not found", *flamingoContext)
	}

	injector, err := app.area.GetInitializedInjector()
	if err != nil {
		return nil, fmtErrorf("app: get initialized injector: %w", err)
	}

	if *dingoInspect {
		inspect(injector)
	}

	if app.eagerSingletons {
		if err := injector.BuildEagerSingletons(false); err != nil {
			return nil, fmtErrorf("app: build eager singletons: %w", err)
		}
	}

	return app, nil
}

// ConfigArea returns the initialized configuration area
func (app *Application) ConfigArea() *config.Area {
	return app.area
}

// App is the default app-runner for flamingo
func App(modules []dingo.Module, options ...ApplicationOption) {
	app, err := NewApplication(modules, options...)
	if err != nil {
		log.Fatal(err)
	}
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

// Run runs the Root Cmd and triggers the standard event
func (app *Application) Run() error {
	injector, err := app.area.GetInitializedInjector()
	if err != nil {
		return fmtErrorf("get initialized injector: %w", err)
	}

	i, err := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo")
	if err != nil {
		return fmtErrorf("app: get flamingo cobra.Command: %w", err)
	}

	rootCmd := i.(*cobra.Command)
	rootCmd.SetArgs(app.flagset.Args())

	i, err = injector.GetInstance(new(eventRouterProvider))
	if err != nil {
		return fmtErrorf("app: get eventRouterProvider: %w", err)
	}
	i.(eventRouterProvider)().Dispatch(context.Background(), new(flamingo.StartupEvent))

	return rootCmd.Execute()
}

func printBinding(of reflect.Type, annotation string, to reflect.Type, provider, instance *reflect.Value, in dingo.Scope) {
	name := of.Name()
	if of.PkgPath() != "" {
		name = of.PkgPath() + "." + name
	}
	if annotation != "" {
		annotation = fmt.Sprintf("[%q]", annotation)
	}
	val := "<unset>"
	if instance != nil {
		val = fmt.Sprintf("value=%q", instance)
	} else if provider != nil {
		val = "provider=" + provider.String()
	} else if to != nil {
		val = "type=" + to.PkgPath() + "." + to.Name()
	}
	scopename := ""
	if in != nil {
		scopename = " (" + reflect.ValueOf(in).String() + ")"
	}
	fmt.Printf("%s%s: %s%s\n", name, annotation, val, scopename)
}

func inspect(injector *dingo.Injector) {
	fmt.Println("Bindings:")
	injector.Inspect(dingo.Inspector{
		InspectBinding: printBinding,
	})

	fmt.Println("\nMultiBindings:")
	injector.Inspect(dingo.Inspector{
		InspectMultiBinding: func(of reflect.Type, index int, annotation string, to reflect.Type, provider, instance *reflect.Value, in dingo.Scope) {
			//fmt.Printf("%d: ", index)
			printBinding(of, annotation, to, provider, instance, in)
		},
	})

	fmt.Println("\nMapBindings:")
	injector.Inspect(dingo.Inspector{
		InspectMapBinding: func(of reflect.Type, key string, annotation string, to reflect.Type, provider, instance *reflect.Value, in dingo.Scope) {
			//fmt.Printf("%s: ", key)
			printBinding(of, annotation, to, provider, instance, in)
		},
	})

	fmt.Println("---")
	injector.Inspect(dingo.Inspector{
		InspectParent: inspect,
	})
}

type servemodule struct {
	router            *web.Router
	server            *http.Server
	eventRouter       flamingo.EventRouter
	logger            flamingo.Logger
	configuredSampler *opencensus.ConfiguredURLPrefixSampler
}

// Inject basic application dependencies
func (a *servemodule) Inject(
	router *web.Router,
	eventRouter flamingo.EventRouter,
	logger flamingo.Logger,
	configuredSampler *opencensus.ConfiguredURLPrefixSampler,
) {
	a.router = router
	a.eventRouter = eventRouter
	a.logger = logger
	a.server = &http.Server{
		Addr: ":3322",
	}
	a.configuredSampler = configuredSampler
}

// Configure dependency injection
func (a *servemodule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).ToInstance(a)

	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

func serveProvider(a *servemodule, logger flamingo.Logger) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", a.server.Addr))
			a.server.Handler = &ochttp.Handler{IsPublicEndpoint: true, Handler: a.router.Handler(), GetStartOptions: a.configuredSampler.GetStartOptions()}

			err := a.listenAndServe()
			if err != nil {
				if err == http.ErrServerClosed {
					logger.Error(err)
				} else {
					logger.Fatal("unexpected error in serving:", err)
				}
			}
		},
	}
	serveCmd.Flags().StringVarP(&a.server.Addr, "addr", "a", ":3322", "addr on which flamingo runs")

	return serveCmd
}

func (a *servemodule) listenAndServe() error {
	a.eventRouter.Dispatch(context.Background(), &flamingo.ServerStartEvent{})
	defer a.eventRouter.Dispatch(context.Background(), &flamingo.ServerShutdownEvent{})

	err := a.server.ListenAndServe()

	return err
}

// Notify upon flamingo Shutdown event
func (a *servemodule) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		a.logger.Info("Shutdown server on ", a.server.Addr)

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Error("unexpected error on server shutdown: ", err)
		}
	}
}
