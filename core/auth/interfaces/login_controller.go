package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/framework/web"
	"flamingo.me/flamingo/v3/framework/web/responder"
	"github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type (
	LoginControllerInterface interface {
		Get(context.Context, *web.Request) web.Response
	}

	// LoginController handles the login redirect
	LoginController struct {
		responder.RedirectAware
		authManager    *application.AuthManager
		parameterHooks []LoginGetParameterHook
	}

	LoginGetParameterHook interface {
		Parameters(context.Context, *web.Request) map[string]string
	}
)

// Inject LoginController dependencies
func (l *LoginController) Inject(
	redirectAware responder.RedirectAware,
	authManager *application.AuthManager,
	ph []LoginGetParameterHook,
) {
	l.RedirectAware = redirectAware
	l.authManager = authManager
	l.parameterHooks = ph
}

// Get handler for logins (redirect)
func (l *LoginController) Get(c context.Context, request *web.Request) web.Response {
	redirecturl, ok := request.Param1("redirecturl")
	if !ok || redirecturl == "" {
		redirecturl = request.Request().Referer()
	}

	if refURL, err := url.Parse(redirecturl); err != nil || refURL.Host != request.Request().Host {
		u, _ := l.authManager.URL(c, "")
		redirecturl = u.String()
	}

	state := uuid.NewV4().String()
	request.Session().Store("auth.state", state)
	request.Session().Store("auth.redirect", redirecturl)

	var parameters []oauth2.AuthCodeOption
	for _, hook := range l.parameterHooks {
		keyValue := hook.Parameters(c, request)
		for key, value := range keyValue {
			parameters = append(parameters, oauth2.SetAuthURLParam(key, value))
		}
	}

	return l.RedirectURL(l.authManager.OAuth2Config(c).AuthCodeURL(state, parameters...))
}
