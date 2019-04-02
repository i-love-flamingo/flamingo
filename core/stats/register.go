package stats

import (
	"expvar"
	"runtime"
	"runtime/debug"

	"go.aoe.com/flamingo/framework/dingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
)

func init() {
	expvar.Publish("routines", expvar.Func(func() interface{} { return runtime.NumGoroutine() }))
}

type (
	// Module registers our profiler
	Module struct {
		RouterRegistry *router.Registry `inject:""`
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
