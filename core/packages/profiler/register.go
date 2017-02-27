package profiler

import (
	"flamingo/core/flamingo/router"
	"flamingo/core/flamingo/service_container"
	"log"
)

type Test struct {
	Router *router.Router `inject:""`
}

func (t Test) PostInject() {
	log.Println(t)
}

func Register(sr *service_container.ServiceContainer) {
	sr.Register(Test{}, "profiler.collector", "logger")
}
