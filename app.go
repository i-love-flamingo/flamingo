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

type option func(config *appconfig)

// ConfigDir configuration option
func ConfigDir(configdir string) func(config *appconfig) {
	return func(config *appconfig) {
		config.configDir = configdir
	}
}

// ChildAreas allows to define additional config areas for roots
func ChildAreas(areas ...*config.Area) func(config *appconfig) {
	return func(config *appconfig) {
		config.childAreas = areas
	}
}

// DefaultContext for flamingo to start with
func DefaultContext(name string) func(config *appconfig) {
	return func(config *appconfig) {
		config.defaultContext = name
	}
}

// SetEagerSingletons controls if eager singletons will be created
func SetEagerSingletons(enabled bool) func(config *appconfig) {
	return func(config *appconfig) {
		config.eagerSingletons = enabled
	}
}

type appconfig struct {
	configDir       string
	childAreas      []*config.Area
	args            []string
	defaultContext  string
	eagerSingletons bool
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

// App is a simple app-runner for flamingo
func App(modules []dingo.Module, options ...option) {
	cfg := &appconfig{
		configDir:      "config",
		args:           os.Args[1:],
		defaultContext: "root",
	}

	for _, option := range options {
		option(cfg)
	}

	fs := flag.NewFlagSet("flamingo", flag.ContinueOnError)
	dingoTraceCircular := fs.Bool("dingo-trace-circular", false, "enable dingo circular tracing")
	flamingoConfigLog := fs.Bool("flamingo-config-log", false, "enable flamingo config logging")
	flamingoConfigCueDebug := fs.String("flamingo-config-cue-debug", "", "query the flamingo cue config loader (use . for root)")
	flamingoContext := fs.String("flamingo-context", cfg.defaultContext, "set flamingo execution context")
	var flamingoConfig arrayFlags
	fs.Var(&flamingoConfig, "flamingo-config", "add additional flamingo yaml config")
	dingoInspect := fs.Bool("dingo-inspect", false, "inspect dingo")

	if err := fs.Parse(cfg.args); err != nil && err != flag.ErrHelp {
		log.Fatal("app: parsing arguments:", err)
	}

	if dingoTraceCircular != nil && *dingoTraceCircular {
		dingo.EnableCircularTracing()
	}

	root := config.NewArea("root", modules, cfg.childAreas...)

	root.Modules = append([]dingo.Module{
		new(framework.InitModule),
		new(zap.Module),
		new(runtime.Module),
		new(cmd.Module),
	}, root.Modules...)
	root.Modules = append(root.Modules, new(appmodule))

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

	if err := config.Load(root, cfg.configDir, configLoadOptions...); err != nil {
		log.Println("app: config load:", err)
		os.Exit(-2)
	}

	areas, err := root.Flat()
	if err != nil {
		log.Fatal("app: flat areas:", err)
	}

	area, ok := areas[*flamingoContext]
	if !ok {
		log.Fatalf("app: context %q not found", *flamingoContext)
	}

	injector, err := area.GetInitializedInjector()
	if err != nil {
		log.Fatal("app: get initialized injector:", err)
	}

	if *dingoInspect {
		inspect(injector)
	}

	if cfg.eagerSingletons {
		if err := injector.BuildEagerSingletons(false); err != nil {
			log.Fatal("app: build eager singletons:", err)
		}
	}

	i, err := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo")
	if err != nil {
		log.Fatal("app: get flamingo cobra.Command:", err)
	}
	rootCmd := i.(*cobra.Command)
	rootCmd.SetArgs(fs.Args())

	i, err = injector.GetInstance(new(eventRouterProvider))
	if err != nil {
		log.Fatal("app: get eventRouterProvider:", err)
	}
	i.(eventRouterProvider)().Dispatch(context.Background(), new(flamingo.StartupEvent))

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("app: rootcmd.Execute:", err)
		os.Exit(-1)
	}
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

type appmodule struct {
	root              *config.Area
	router            *web.Router
	server            *http.Server
	eventRouter       flamingo.EventRouter
	logger            flamingo.Logger
	configuredSampler *opencensus.ConfiguredURLPrefixSampler
}

// Inject basic application dependencies
func (a *appmodule) Inject(
	root *config.Area,
	router *web.Router,
	eventRouter flamingo.EventRouter,
	logger flamingo.Logger,
	configuredSampler *opencensus.ConfiguredURLPrefixSampler,
) {
	a.root = root
	a.router = router
	a.eventRouter = eventRouter
	a.logger = logger
	a.server = &http.Server{
		Addr: ":3322",
	}
	a.configuredSampler = configuredSampler
}

// Configure dependency injection
func (a *appmodule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).ToInstance(a)

	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

func serveProvider(a *appmodule, logger flamingo.Logger) *cobra.Command {
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

func (a *appmodule) listenAndServe() error {
	a.eventRouter.Dispatch(context.Background(), &flamingo.ServerStartEvent{})
	defer a.eventRouter.Dispatch(context.Background(), &flamingo.ServerShutdownEvent{})

	err := a.server.ListenAndServe()

	return err
}

// Notify upon flamingo Shutdown event
func (a *appmodule) Notify(ctx context.Context, event flamingo.Event) {
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
