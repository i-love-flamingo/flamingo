package requestTask

import (
	"flamingo.me/flamingo/v3/framework/dingo"
	"flamingo.me/flamingo/v3/framework/router"
)

type Module struct{}

func (m *Module) Configure(injector *dingo.Injector) {
	injector.BindMulti(new(router.Filter)).To(new(filter))
}
