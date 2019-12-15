package cache

import (
	"flamingo.me/dingo"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

// Module for gotemplate engine
type Module struct {
	cacheConfig          config.Map
	cacheFrontendFactory *HTTPFrontendFactory
	logger               flamingo.Logger
	override             bool
}

//Inject for modul
func (m *Module) Inject(cacheFrontendFactory *HTTPFrontendFactory, logger flamingo.Logger, config *struct {
	CacheConfig config.Map `inject:"config:core.cache.httpFrontendFactory,optional"`
	Override    bool       `inject:"config:core.cache.overrideBindings,optional"`
}) {
	m.cacheFrontendFactory = cacheFrontendFactory
	m.logger = logger.WithField(flamingo.LogKeyModule, "cache").WithField(flamingo.LogKeyCategory, "module")
	if config != nil {
		m.cacheConfig = config.CacheConfig
		m.override = config.Override
	}
}

// Configure DI
func (m *Module) Configure(injector *dingo.Injector) {
	for k := range m.cacheConfig {
		cache, err := m.cacheFrontendFactory.BuildConfiguredCache(k)
		if err != nil {
			m.logger.Fatal(err)
		} else {
			if m.override {
				injector.Override((*HTTPFrontend)(nil), k).ToInstance(cache)
			} else {
				injector.Bind((*HTTPFrontend)(nil)).AnnotatedWith(k).ToInstance(cache)
			}
		}
	}
}
