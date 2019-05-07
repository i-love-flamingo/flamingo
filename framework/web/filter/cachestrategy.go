package filter

import (
	"context"
	"net/http"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/web"
)

// DefaultCacheStrategyModule is a flamingo module to set up a web filter which injects a default cache strategy
type DefaultCacheStrategyModule struct{}

// Configure the Module
func (m *DefaultCacheStrategyModule) Configure(injector *dingo.Injector) {
	injector.BindMulti((*web.Filter)(nil)).To(cacheStrategyFilter{})
}

type cacheStrategyFilter struct {
	DefaultIsReuseable             bool    `inject:"config:flamingo.web.filter.cachestrategy.default.isReusable,optional"`
	DefaultRevalidateEachTime      bool    `inject:"config:flamingo.web.filter.cachestrategy.default.revalidateEachTime,optional"`
	DefaultMaxCacheLifetime        float64 `inject:"config:flamingo.web.filter.cachestrategy.default.maxCacheLifetime,optional"`
	DefaultAllowIntermediateCaches bool    `inject:"config:flamingo.web.filter.cachestrategy.default.allowIntermediateCaches,optional"`
}

// Filter sets the cache strategy for responses
func (f *cacheStrategyFilter) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	response := chain.Next(ctx, r, w)
	if r.Request().Method != http.MethodGet {
		return response
	}

	switch response := response.(type) {
	case *web.RenderResponse:
		f.setDefault(&response.Response)
	case *web.DataResponse:
		f.setDefault(&response.Response)
	}

	return response
}

func (f *cacheStrategyFilter) setDefault(response *web.Response) {
	if response.CacheDirective != nil {
		return
	}

	response.CacheDirective = web.CacheDirectiveBuilder{
		IsReusable:              f.DefaultIsReuseable,
		RevalidateEachTime:      f.DefaultRevalidateEachTime,
		AllowIntermediateCaches: f.DefaultAllowIntermediateCaches,
		MaxCacheLifetime:        int(f.DefaultMaxCacheLifetime),
	}.Build()
}
