package filter

import (
	"context"
	"net/http"

	"flamingo.me/flamingo/v3/framework/web"
)

type (
	//CacheStrategy - a filter that sets CacheDirective if not present
	cacheStrategy struct {
		DefaultIsReuseable             bool    `inject:"config:flamingo.web.filter.cachestrategy.default.isReusable,optional"`
		DefaultRevalidateEachTime      bool    `inject:"config:flamingo.web.filter.cachestrategy.default.revalidateEachTime,optional"`
		DefaultMaxCacheLifetime        float64 `inject:"config:flamingo.web.filter.cachestrategy.default.maxCacheLifetime,optional"`
		DefaultAllowIntermediateCaches bool    `inject:"config:flamingo.web.filter.cachestrategy.default.allowIntermediateCaches,optional"`
	}
)

//Filter - implements flamingo filter interface
func (f *cacheStrategy) Filter(ctx context.Context, r *web.Request, w http.ResponseWriter, chain *web.FilterChain) web.Result {
	response := chain.Next(ctx, r, w)
	if r.Request().Method != http.MethodGet {
		return response
	}
	switch typedResponse := response.(type) {
	case *web.RenderResponse:
		f.setDefault(&typedResponse.Response)
	}
	return response
}

//setDefault - sets default on Basic response
func (f *cacheStrategy) setDefault(response *web.Response) {
	if response.CacheDirective != nil {
		return
	}
	response.CacheDirective = web.NewCacheDirectiveBuilder().SetIsReusable(f.DefaultIsReuseable).
		SetRevalidateEachTime(f.DefaultRevalidateEachTime).
		SetAllowIntermediateCaches(f.DefaultAllowIntermediateCaches).
		SetMaxCacheLifetime(int(f.DefaultMaxCacheLifetime)).
		Build()
}
