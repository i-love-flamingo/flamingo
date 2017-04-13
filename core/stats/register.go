package stats

import (
	"expvar"
	"flamingo/framework/dependencyinjection"
	"flamingo/framework/router"
	"flamingo/framework/web"
	"runtime"
	"runtime/debug"
)

func init() {
	expvar.Publish("routines", expvar.Func(func() interface{} { return runtime.NumGoroutine() }))
}

func Register(c *dependencyinjection.Container) {
	c.Register(func(r *router.Router) {
		r.Handle("stats.index", func(c web.Context) interface{} {
			return map[string]int{
				"routines": runtime.NumGoroutine(),
			}
		})
		r.Handle("stats.expvar", expvar.Handler())
		r.Handle("stats.gc", func(c web.Context) interface{} {
			runtime.GC()
			debug.FreeOSMemory()
			return nil
		})
		r.Route("/stats", "stats.index")
		r.Route("/stats/expvar", "stats.expvar")
		r.Route("/stats/gc", "stats.gc")
	}, router.RouterRegister)
}
