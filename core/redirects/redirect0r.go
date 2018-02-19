package redirects

import (
	"errors"
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
		responder.RedirectAware      `inject:""`
		responder.ErrorAware         `inject:""`
		cache.Backend                `inject:""`
		Logger       flamingo.Logger `inject:""`
		redirectData []infrastructure.CsvContent
	}
)

var redirectData []infrastructure.CsvContent

func init() {
	redirectData = infrastructure.GetRedirectData()
}

func (r *redirect0r) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	contextPath := ctx.Request().RequestURI

	for i := range redirectData {
		originalPath := redirectData[i].OriginalPath

		if originalPath == contextPath {
			redirectTarget := redirectData[i].RedirectTarget
			httpCode := redirectData[i].HttpStatusCode

			r.Logger.Debugf("Redirecting from %s to %s by %d", originalPath, redirectTarget, httpCode)

			switch code := httpCode; code {
			case http.StatusMovedPermanently:
				return r.RedirectPermanentURL(fmt.Sprintf("%s", redirectTarget))
			case http.StatusFound:
				return r.RedirectURL(fmt.Sprintf("%s", redirectTarget))
			}

			return r.ErrorAware.ErrorNotFound(ctx, errors.New("page not found"))
		}
	}

	return chain.Next(ctx, w)
}
