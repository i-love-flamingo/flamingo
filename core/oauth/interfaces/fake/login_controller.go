package fake

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// LoginController fake implementation
	LoginController struct {
		responder     *web.Responder
		authManager   *application.AuthManager
		loginTemplate string
		router        web.ReverseRouter
	}
)

// Inject dependencies
func (l *LoginController) Inject(
	responder *web.Responder,
	authManager *application.AuthManager,
	cfg *struct {
		FakeLoginTemplate string `inject:"config:auth.fakeLoginTemplate"`
	},
	router web.ReverseRouter,
) {
	l.responder = responder
	l.authManager = authManager
	l.loginTemplate = cfg.FakeLoginTemplate
	l.router = router
}

// Get http action
func (l *LoginController) Get(ctx context.Context, request *web.Request) web.Result {
	redirectURL, ok := request.Params["redirecturl"]
	if !ok || redirectURL == "" {
		redirectURL = request.Request().Referer()
	}

	if refURL, err := url.Parse(redirectURL); err != nil || refURL.Host != request.Request().Host {
		u, _ := l.router.Absolute(request, "", nil)
		redirectURL = u.String()
	}

	if redirectURL != "" {
		request.Session().Store("auth.redirect", redirectURL)
	}

	if l.loginTemplate != "" {
		return l.responder.Render(l.loginTemplate, nil)
	}

	return l.responder.RouteRedirect("auth.callback", nil)
}
