package requestTask

import (
	"flamingo.me/flamingo/framework/dingo"
	"flamingo.me/flamingo/framework/router"
)

type Module struct{}

func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(router.Filter)).To(new(filter))
}
