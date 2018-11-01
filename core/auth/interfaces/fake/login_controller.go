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
		responder.RenderAware

		authManager *application.AuthManager

		loginTemplate string
	}
)

func (l *LoginController) Inject(
	redirectAware responder.RedirectAware,
	renderAware responder.RenderAware,
	authManager *application.AuthManager,
	cfg *struct {
		FakeLoginTemplate string `inject:"config:auth.fakeLoginTemplate"`
	},
) {
	l.RedirectAware = redirectAware
	l.RenderAware = renderAware
	l.authManager = authManager
	l.loginTemplate = cfg.FakeLoginTemplate
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
		request.Session().Store("auth.redirect", redirectUrl)
	}

	if l.loginTemplate != "" {
		return l.Render(ctx, l.loginTemplate, nil)
	}

	return l.Redirect("auth.callback", nil)
}
