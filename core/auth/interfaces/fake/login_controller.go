package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	LoginController struct {
		responder.RedirectAware
		authManager *application.AuthManager
	}
)

func (l *LoginController) Inject(
	redirectAware responder.RedirectAware,
	authManager *application.AuthManager,
) {
	l.RedirectAware = redirectAware
	l.authManager = authManager
}

func (l *LoginController) Get(ctx context.Context, request *web.Request) web.Response {
	redirectUrl, ok := request.Param1("redirecturl")
	if !ok || redirectUrl == "" {
		redirectUrl = request.Request().Referer()
	}

	if refURL, err := url.Parse(redirectUrl); err != nil || refURL.Host != request.Request().Host {
		u, _ := l.authManager.URL(ctx, "")
		redirectUrl = u.String()
	}

	if redirectUrl != "" {
		request.Session().Values["auth.redirect"] = redirectUrl
	}

	return l.Redirect("auth.callback", nil)
}
