package stats

import (
	"expvar"
	"flamingo/framework/dingo"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"runtime"
	"runtime/debug"
)

func init() {
	expvar.Publish("routines", expvar.Func(func() interface{} { return runtime.NumGoroutine() }))
}

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.RouterRegistry `inject:""`
	}
)

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	m.RouterRegistry.Handle("stats.index", func(c web.Context) interface{} {
		return map[string]int{
			"routines": runtime.NumGoroutine(),
		}
	})
	m.RouterRegistry.Handle("stats.gc", func(c web.Context) interface{} {
		runtime.GC()
		debug.FreeOSMemory()
		return nil
	})
	m.RouterRegistry.Handle("stats.expvar", expvar.Handler())
	m.RouterRegistry.Route("/stats", "stats.index")
	m.RouterRegistry.Route("/stats/expvar", "stats.expvar")
	m.RouterRegistry.Route("/stats/gc", "stats.gc")
}
