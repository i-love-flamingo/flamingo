package redirects

import (
	"fmt"
	"net/http"

	"go.aoe.com/flamingo/core/redirects/infrastructure"
	"go.aoe.com/flamingo/framework/flamingo"
	"go.aoe.com/flamingo/framework/router"
	"go.aoe.com/flamingo/framework/web"
	"go.aoe.com/flamingo/framework/web/responder"

	"go.aoe.com/flamingo/core/cache"
)

type (
	redirect0r struct {
		responder.RedirectAware                `inject:""`
		cache.Backend                          `inject:""`
		Logger flamingo.Logger                 `inject:""`
	}
)

var redirectData []infrastructure.CsvContent

func init() {
	redirectData = infrastructure.GetRedirectData()
}

func (r *redirect0r) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	for i := range redirectData {
		if redirectData[i].OriginalPath == ctx.Request().URL.RawPath {
			r.redirectFlow(ctx.Request(), w, redirectData[i])
		}
	}

	return chain.Next(ctx, w)
}

func (r *redirect0r) redirectFlow(req *http.Request, w http.ResponseWriter, redirectInfo infrastructure.CsvContent) {
	switch code := redirectInfo.HttpStatusCode; code {
	case http.StatusGone:
		r.Redirect(fmt.Sprintf("/%d", http.StatusNotFound), nil)
	case http.StatusMovedPermanently:
		r.RedirectPermanent(fmt.Sprintf("/%s", redirectInfo.RedirectTarget), nil)
	case http.StatusFound:
		r.Redirect(fmt.Sprintf("/%s", redirectInfo.RedirectTarget), nil)
	}

	r.Redirect(fmt.Sprintf("/%d", http.StatusNotFound), nil)
}
