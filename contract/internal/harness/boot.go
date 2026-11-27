// Package harness starts a Flamingo app inside the test process.
package harness

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3"
	flamingozap "flamingo.me/flamingo/v3/core/zap"
	"flamingo.me/flamingo/v3/framework/config"
	flamingoFramework "flamingo.me/flamingo/v3/framework/flamingo"
	flamingoHTTP "flamingo.me/flamingo/v3/framework/http"
)

const (
	bootTimeout    = 30 * time.Second
	dirPermission  = 0o750
	filePermission = 0o600
)

var errAppPanic = errors.New("app run: panic")

// Config describes the app that Boot starts.
type Config struct {
	Files map[string]string // files of the config dir, path relative to it: "config.yml", "en/config.yml", "config.cue"
	Areas []*config.Area    // child areas of root, given to flamingo.ChildAreas
	Env   map[string]string // set with t.Setenv before the config loads

	Options []flamingo.ApplicationOption // more options for flamingo.NewApplication, applied after the ones of Boot
}

// App is a running app.
type App struct {
	URL    string       // "http://127.0.0.1:<port>"
	Client *http.Client // cookie jar, does not follow redirects

	area *config.Area // Application.ConfigArea()
	logs *observer.ObservedLogs
}

// portNotifier sends the port of the first started server to a channel.
type portNotifier struct {
	ports chan string
}

// Notify implements the flamingo event subscriber.
func (n portNotifier) Notify(_ context.Context, event flamingoFramework.Event) {
	if started, ok := event.(*flamingoFramework.ServerStartEvent); ok {
		select {
		case n.ports <- started.Port:
		default:
		}
	}
}

// Isolate gives dingo new Singleton and ChildSingleton scopes and clears CONTEXT and CONTEXTFILE,
// so an app sees neither instances of an earlier app nor the config context of the shell.
// It uses t.Setenv, so a test that calls it cannot be parallel.
func Isolate(tb testing.TB) {
	tb.Helper()

	dingo.Singleton = dingo.NewSingletonScope()
	dingo.ChildSingleton = dingo.NewChildSingletonScope()

	tb.Setenv("CONTEXT", "")
	tb.Setenv("CONTEXTFILE", "")
}

// Boot starts an app on a free port and stops it when the test ends. The test cannot be parallel.
// Only one app can be in use at a time: stopping one interrupts the test process, which stops all of them.
func Boot(tb testing.TB, modules []dingo.Module, cfg Config) *App {
	tb.Helper()

	app, err := TryBoot(tb, modules, cfg)
	if err != nil {
		tb.Fatalf("boot: %v", err)
	}

	return app
}

// TryBoot is Boot, but it returns the error of flamingo.NewApplication, such as a config or injection error.
func TryBoot(tb testing.TB, modules []dingo.Module, cfg Config) (*App, error) {
	tb.Helper()

	Isolate(tb)

	for name, value := range cfg.Env {
		tb.Setenv(name, value)
	}

	dir := writeFiles(tb, cfg.Files)

	notifier := portNotifier{ports: make(chan string, 1)}

	modules = append(slices.Clip(modules), dingo.ModuleFunc(func(injector *dingo.Injector) {
		flamingoFramework.BindEventSubscriber(injector).ToInstance(notifier)
	}))

	// This logger takes the place of zap.Module and records the entries for App.Logs.
	core, logs := observer.New(zap.DebugLevel)
	logger := flamingozap.NewLogger(zap.New(core), flamingozap.WithArea("root"))
	loggerModule := dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.Bind(new(flamingoFramework.Logger)).ToInstance(logger)
	})

	options := append([]flamingo.ApplicationOption{
		flamingo.WithCustomLogger(loggerModule),
		flamingo.ConfigDir(dir),
		flamingo.ChildAreas(cfg.Areas...),
		flamingo.WithArgs(
			"--flamingo-config", `flamingo.systemendpoint.serviceAddr: "127.0.0.1:0"`,
			"serve", "--addr", "127.0.0.1:0",
		),
	}, cfg.Options...)

	app, err := flamingo.NewApplication(modules, options...)
	if err != nil {
		return nil, fmt.Errorf("new application: %w", err)
	}

	done := make(chan error, 1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				done <- fmt.Errorf("%w: %v\n%s", errAppPanic, r, debug.Stack())
			}
		}()

		done <- app.Run()
	}()

	var port string

	select {
	case port = <-notifier.ports:
	case err := <-done:
		tb.Fatalf("boot: app stopped before the server started: %v", err)
	case <-time.After(bootTimeout):
		tb.Fatalf("boot: no server start after %s", bootTimeout)
	}

	tb.Cleanup(func() {
		stop(tb, done)
	})

	jar, err := cookiejar.New(nil)
	if err != nil {
		tb.Fatalf("boot: cookie jar: %v", err)
	}

	return &App{
		URL: "http://127.0.0.1:" + port,
		Client: &http.Client{
			Jar: jar,
			CheckRedirect: func(*http.Request, []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		area: app.ConfigArea(),
		logs: logs,
	}, nil
}

func writeFiles(tb testing.TB, files map[string]string) string {
	tb.Helper()

	dir := tb.TempDir()

	for name, content := range files {
		path := filepath.Join(dir, filepath.FromSlash(name))

		if err := os.MkdirAll(filepath.Dir(path), dirPermission); err != nil {
			tb.Fatalf("boot: create dir for %s: %v", name, err)
		}

		if err := os.WriteFile(path, []byte(content), filePermission); err != nil {
			tb.Fatalf("boot: write %s: %v", name, err)
		}
	}

	return dir
}

// stop interrupts the process, as a user does, and waits for app.Run to return.
func stop(tb testing.TB, done <-chan error) {
	tb.Helper()

	select {
	case err := <-done:
		if err != nil {
			tb.Errorf("shutdown: app run ended before the test end: %v", err)
		}

		return
	default:
	}

	self, err := os.FindProcess(os.Getpid())
	if err != nil {
		tb.Errorf("shutdown: find own process: %v", err)

		return
	}

	if err := self.Signal(os.Interrupt); err != nil {
		tb.Errorf("shutdown: send interrupt: %v", err)

		return
	}

	select {
	case err := <-done:
		if err != nil {
			tb.Errorf("shutdown: app run: %v", err)
		}
	case <-time.After(bootTimeout):
		tb.Errorf("shutdown: app still runs after %s", bootTimeout)
	}
}

// Injector returns the initialized injector of the app area.
func (a *App) Injector(tb testing.TB) *dingo.Injector {
	tb.Helper()

	injector, err := a.area.GetInitializedInjector()
	if err != nil {
		tb.Fatalf("injector of area %s: %v", a.area.Name, err)
	}

	return injector
}

// Response is the answer to App.Get, with the body read.
type Response struct {
	Status int         // HTTP status code
	Header http.Header // response header
	Body   string      // whole response body
}

// Get sends a GET request to the app, reads and closes the body.
// The optional header values are name and value pairs. A "Host" pair sets Request.Host,
// because the client does not send a "Host" entry of Request.Header.
func (a *App) Get(tb testing.TB, path string, header ...string) Response {
	tb.Helper()

	if len(header)%2 != 0 {
		tb.Fatalf("get %s: header needs name and value pairs, got %d values", path, len(header))
	}

	req, err := http.NewRequestWithContext(tb.Context(), http.MethodGet, a.URL+path, nil)
	if err != nil {
		tb.Fatalf("get %s: %v", path, err)
	}

	for i := 0; i < len(header); i += 2 {
		if http.CanonicalHeaderKey(header[i]) == "Host" {
			req.Host = header[i+1]

			continue
		}

		req.Header.Set(header[i], header[i+1])
	}

	resp, err := a.Client.Do(req)
	if err != nil {
		tb.Fatalf("get %s: %v", path, err)
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		tb.Fatalf("get %s: read body: %v", path, err)
	}

	return Response{Status: resp.StatusCode, Header: resp.Header, Body: string(body)}
}

// Logs returns all log entries of the app.
func (a *App) Logs() []observer.LoggedEntry {
	return a.logs.All()
}

// PassThroughWrapper binds a handler wrapper that does nothing. prefixrouter needs a wrapper, and opencensus
// binds its own only in the first app of a process, so do not use this together with opencensus.Module.
func PassThroughWrapper() dingo.Module {
	return dingo.ModuleFunc(func(injector *dingo.Injector) {
		injector.Bind(new(flamingoHTTP.HandlerWrapper)).ToProvider(func() flamingoHTTP.HandlerWrapper {
			return func(handler http.Handler) http.Handler { return handler }
		})
	})
}
