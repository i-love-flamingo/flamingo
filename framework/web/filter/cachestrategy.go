package filter

import (
	"context"
	"flamingo.me/flamingo/v3/framework/web"
	"net/http"
)

type (
	//CacheStrategy - a filter that sets CacheDirective if not present
	CacheStrategy struct {
		DefaultIsReuseable             bool    `inject:"config:flamingo.web.filter.cachestrategy.default.isReusable,optional"`
		DefaultRevalidateEachTime      bool    `inject:"config:flamingo.web.filter.cachestrategy.default.revalidateEachTime,optional"`
		DefaultMaxCacheLifetime        float64 `inject:"config:flamingo.web.filter.cachestrategy.default.maxCacheLifetime,optional"`
		DefaultAllowIntermediateCaches bool    `inject:"config:flamingo.web.filter.cachestrategy.default.allowIntermediateCaches,optional"`
	}
)

//Filter - implements flamingo filter interface
func (f *CacheStrategy) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	response := chain.Next(ctx, r, w)
	if defaultResponse, ok := response.(*web.Response); ok {
		if defaultResponse.CacheDirectives == nil && r.Request().Method == http.MethodGet {
			cacheStrategy := web.CacheStrategy{}
			cacheStrategy.SetIsReusable(f.DefaultIsReuseable)
			cacheStrategy.SetRevalidateEachTime(f.DefaultRevalidateEachTime)
			cacheStrategy.SetAllowIntermediateCaches(f.DefaultAllowIntermediateCaches)
			cacheStrategy.SetMaxCacheLifetime(int(f.DefaultMaxCacheLifetime))
			defaultResponse.CacheDirectives = cacheStrategy.Build()
			return defaultResponse
		}
	}
	return response
}
