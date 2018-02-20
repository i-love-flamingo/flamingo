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
)

type (
	redirector struct {
		responder.RedirectAware `inject:""`
		responder.ErrorAware    `inject:""`
		Logger flamingo.Logger  `inject:""`
	}
)

var redirectDataMap map[string]infrastructure.CsvContent

func init() {
	redirectData := infrastructure.GetRedirectData()

	redirectDataMap = make(map[string]infrastructure.CsvContent)

	for i := range redirectData {
		redirectDataMap[redirectData[i].OriginalPath] = redirectData[i]
	}
}

func (r *redirector) Filter(ctx web.Context, w http.ResponseWriter, chain *router.FilterChain) web.Response {
	contextPath := ctx.Request().RequestURI

	if currentRedirect, ok := redirectDataMap[contextPath]; ok {
		httpCode := currentRedirect.HttpStatusCode
		originalPath := currentRedirect.OriginalPath
		redirectTarget := currentRedirect.RedirectTarget

		r.Logger.Debugf("Redirecting from %s to %s by %d", originalPath, redirectTarget, httpCode)

		switch code := httpCode; code {
		case http.StatusMovedPermanently:
			return r.RedirectPermanentURL(redirectTarget)
		case http.StatusFound:
			return r.RedirectURL(fmt.Sprintf("%s", redirectTarget))
		case http.StatusGone:
			return r.ErrorAware.ErrorWithCode(ctx, errors.New("page is gone"), http.StatusGone)
		}

		return r.ErrorAware.ErrorNotFound(ctx, errors.New("page not found"))
	}

	return chain.Next(ctx, w)
}
