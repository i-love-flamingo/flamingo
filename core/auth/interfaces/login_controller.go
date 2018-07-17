package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
	"github.com/satori/go.uuid"
)

type (
	// LoginController handles the login redirect
	LoginController struct {
		responder.RedirectAware
		authManager *application.AuthManager
		myHost      string
	}
)

// Inject LoginController dependencies
func (l *LoginController) Inject(
	redirectAware responder.RedirectAware,
	authManager *application.AuthManager,
	config *struct {
		MyHost string `inject:"config:auth.myhost"`
	},
) {
	l.RedirectAware = redirectAware
	l.authManager = authManager
	l.myHost = config.MyHost
}

// Get handler for logins (redirect)
func (l *LoginController) Get(_ context.Context, request *web.Request) web.Response {
	redirecturl, ok := request.Param1("redirecturl")
	if !ok || redirecturl == "" {
		redirecturl = request.Request().Referer()
	}

	if refURL, err := url.Parse(redirecturl); err != nil || refURL.Host != request.Request().Host {
		redirecturl = l.myHost
	}

	state := uuid.NewV4().String()
	request.Session().Values["auth.state"] = state
	request.Session().Values["auth.redirect"] = redirecturl

	return l.RedirectURL(l.authManager.OAuth2Config().AuthCodeURL(state))
}
