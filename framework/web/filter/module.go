package filter

import (
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//CacheStrategyModule - Dingo module that registers the cacheStrategy filter
	CacheStrategyModule struct {
	}
)

// Configure the InitModule
func (m *CacheStrategyModule) Configure(injector *dingo.Injector) {
	injector.BindMulti((*web.Filter)(nil)).To(cacheStrategy{})
}
