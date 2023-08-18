package flamingo

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"flamingo.me/dingo"
	"github.com/spf13/cobra"
	"go.opencensus.io/plugin/ochttp"

	"flamingo.me/flamingo/v3/core/runtime"
	"flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework"
	"flamingo.me/flamingo/v3/framework/cmd"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/opencensus"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// Application contains a main flamingo application
	Application struct {
		configDir       string
		childAreas      []*config.Area
		area            *config.Area
		args            []string
		routesModules   []web.RoutesModule
		loggerModule    dingo.Module
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

// WithRoutes configures a given RoutesModule for usage in the flamingo app
func WithRoutes(routesModule web.RoutesModule) ApplicationOption {
	return func(config *Application) {
		config.routesModules = append(config.routesModules, routesModule)
	}
}

// WithCustomLogger allows to use custom logger modules for flamingo app, if nothing available default will be used
func WithCustomLogger(logger dingo.Module) ApplicationOption {
	return func(config *Application) {
		config.loggerModule = logger
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
		loggerModule:   new(zap.Module),
	}

	for _, option := range options {
		option(app)
	}

	var flamingoConfig arrayFlags

	app.flagset = flag.NewFlagSet("flamingo", flag.ContinueOnError)
	dingoTraceCircular := app.flagset.Bool("dingo-trace-circular", false, "enable dingo circular tracing")
	flamingoConfigLog := app.flagset.Bool("flamingo-config-log", false, "enable flamingo config logging")
	flamingoConfigCueDebug := app.flagset.String("flamingo-config-cue-debug", "", "query the flamingo cue config loader (use . for root)")
	flamingoContext := app.flagset.String("flamingo-context", app.defaultContext, "set flamingo execution context")
	app.flagset.Var(&flamingoConfig, "flamingo-config", "add additional flamingo yaml config")
	dingoInspect := app.flagset.Bool("dingo-inspect", false, "inspect dingo")

	if err := app.flagset.Parse(app.args); err != nil && err != flag.ErrHelp {
		return nil, fmt.Errorf("app: parsing arguments: %w", err)
	}

	if dingoTraceCircular != nil && *dingoTraceCircular {
		dingo.EnableCircularTracing()
	}

	modules = append([]dingo.Module{
		new(framework.InitModule),
		app.loggerModule,
		new(runtime.Module),
		new(cmd.Module),
	}, modules...)
	modules = append(modules, new(servemodule))
	for _, routesModule := range app.routesModules {
		modules = append(modules, dingo.ModuleFunc(func(injector *dingo.Injector) {
			web.BindRoutes(injector, routesModule)
		}))
	}

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
		return nil, fmt.Errorf("app: config load: %w", err)
	}

	areas, err := root.Flat()
	if err != nil {
		return nil, fmt.Errorf("app: flat areas: %w", err)
	}

	var ok bool
	app.area, ok = areas[*flamingoContext]
	if !ok {
		return nil, fmt.Errorf("app: context %q not found", *flamingoContext)
	}

	injector, err := app.area.GetInitializedInjector()
	if err != nil {
		return nil, fmt.Errorf("app: get initialized injector: %w", err)
	}

	if *dingoInspect {
		inspect(injector)
	}

	if app.eagerSingletons {
		if err := injector.BuildEagerSingletons(false); err != nil {
			return nil, fmt.Errorf("app: build eager singletons: %w", err)
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
		return fmt.Errorf("get initialized injector: %w", err)
	}

	instance, err := injector.GetAnnotatedInstance(new(cobra.Command), "flamingo")
	if err != nil {
		return fmt.Errorf("app: get flamingo cobra.Command: %w", err)
	}

	rootCmd := instance.(*cobra.Command)
	rootCmd.SetArgs(app.flagset.Args())

	instance, err = injector.GetInstance(new(eventRouterProvider))
	if err != nil {
		return fmt.Errorf("app: get eventRouterProvider: %w", err)
	}
	instance.(eventRouterProvider)().Dispatch(context.Background(), new(flamingo.StartupEvent))

	return rootCmd.Execute()
}

func typeName(target reflect.Type) string {
	var name string

	for target.Kind() == reflect.Ptr {
		target = target.Elem()
	}

	if target.Kind() == reflect.Slice {
		name += "[]"
		target = target.Elem()
	}

	if target.Kind() == reflect.Ptr {
		name += "*"
		target = target.Elem()
	}

	if target.PkgPath() != "" {
		name += target.PkgPath() + "."
	}

	name += target.Name()

	return name
}

const truncMax = 25

func trunc(s string) string {
	if len(s) > truncMax {
		return s[:truncMax] + "..."
	}
	return s
}

func printBinding(of reflect.Type, annotation string, to reflect.Type, provider, instance *reflect.Value, in dingo.Scope) {
	name := typeName(of)
	if annotation != "" {
		annotation = fmt.Sprintf("(%q)", annotation)
	}
	val := "<unset>"
	if instance != nil {
		val = trunc(fmt.Sprintf("%v", instance.Interface()))
	} else if provider != nil {
		val = "provider=" + provider.String()
	} else if to != nil {
		val = "type=" + typeName(to)
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
			// fmt.Printf("%d: ", index)
			printBinding(of, annotation, to, provider, instance, in)
		},
	})

	fmt.Println("\nMapBindings:")
	injector.Inspect(dingo.Inspector{
		InspectMapBinding: func(of reflect.Type, key string, annotation string, to reflect.Type, provider, instance *reflect.Value, in dingo.Scope) {
			// fmt.Printf("%s: ", key)
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
	certFile, keyFile string
	publicEndpoint    bool
}

// Inject basic application dependencies
func (a *servemodule) Inject(
	router *web.Router,
	eventRouter flamingo.EventRouter,
	logger flamingo.Logger,
	configuredSampler *opencensus.ConfiguredURLPrefixSampler,
	cfg *struct {
		Port           int  `inject:"config:core.serve.port"`
		PublicEndpoint bool `inject:"config:flamingo.opencensus.publicEndpoint,optional"`
	},
) {
	a.router = router
	a.eventRouter = eventRouter
	a.logger = logger
	a.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		ReadHeaderTimeout: 10 * time.Second,
	}
	a.configuredSampler = configuredSampler
	a.publicEndpoint = cfg.PublicEndpoint
}

// Configure dependency injection
func (a *servemodule) Configure(injector *dingo.Injector) {
	flamingo.BindEventSubscriber(injector).ToInstance(a)

	injector.BindMulti(new(cobra.Command)).ToProvider(func() *cobra.Command {
		return serveProvider(a, a.logger)
	})
}

// CueConfig for the module
func (a *servemodule) CueConfig() string {
	return `core: serve: port: >= 0 & <= 65535 | *3322`
}

func serveProvider(serveModule *servemodule, logger flamingo.Logger) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Default serve command - starts on Port 3322",
		Run: func(cmd *cobra.Command, args []string) {
			serveModule.server.Handler = &ochttp.Handler{IsPublicEndpoint: serveModule.publicEndpoint, Handler: serveModule.router.Handler(), GetStartOptions: serveModule.configuredSampler.GetStartOptions()}

			err := serveModule.listenAndServe()
			if err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					logger.Info(err)
				} else {
					logger.Fatal("unexpected error in serving:", err)
				}
			}
		},
	}
	serveCmd.Flags().StringVarP(&serveModule.server.Addr, "addr", "a", serveModule.server.Addr, "addr on which flamingo runs")
	serveCmd.Flags().StringVarP(&serveModule.certFile, "certFile", "c", "", "certFile to enable HTTPS")
	serveCmd.Flags().StringVarP(&serveModule.keyFile, "keyFile", "k", "", "keyFile to enable HTTPS")

	return serveCmd
}

func (a *servemodule) listenAndServe() error {
	listener, err := net.Listen("tcp", a.server.Addr)
	if err != nil {
		return err
	}

	addr := listener.Addr().String()
	a.logger.Info(fmt.Sprintf("Starting HTTP Server at %s .....", addr))

	port := addr[strings.LastIndex(addr, ":")+1:]
	a.eventRouter.Dispatch(context.Background(), &flamingo.ServerStartEvent{Port: port})
	defer a.eventRouter.Dispatch(context.Background(), &flamingo.ServerShutdownEvent{})

	if a.certFile != "" && a.keyFile != "" {
		return a.server.ServeTLS(listener, a.certFile, a.keyFile)
	}

	return a.server.Serve(listener)
}

// Notify upon flamingo Shutdown event
func (a *servemodule) Notify(ctx context.Context, event flamingo.Event) {
	if _, ok := event.(*flamingo.ShutdownEvent); ok {
		if a.server.Handler == nil {
			// server not running, nothing to shut down
			return
		}
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		a.logger.Info("Shutdown server on ", a.server.Addr)

		err := a.server.Shutdown(ctx)
		if err != nil {
			a.logger.Error("unexpected error on server shutdown: ", err)
		}
	}
}
