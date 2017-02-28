package profiler

import "flamingo/core/flamingo/service_container"

func Register(sc *service_container.ServiceContainer) {
	sc.Route("/_profiler/view/{Profile}", "_profiler.view")
	sc.Handle("_profiler.view", new(ProfileController))

	sc.Register(func() *DefaultProfiler { return new(DefaultProfiler) }, "event.subscriber")
}
