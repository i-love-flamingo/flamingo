package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/oauth/application"
	"flamingo.me/flamingo/v3/framework/web"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
)

type (
	// LoginControllerInterface is the callback HTTP action provider
	LoginControllerInterface interface {
		Get(context.Context, *web.Request) web.Result
	}

	// LoginController handles the login redirect
	LoginController struct {
		responder      *web.Responder
		authManager    *application.AuthManager
		parameterHooks []LoginGetParameterHook
		router         web.ReverseRouter
	}

	// LoginGetParameterHook helper to inject additional GET parameters for logins
	LoginGetParameterHook interface {
		Parameters(context.Context, *web.Request) map[string]string
	}
)

// Inject LoginController dependencies
func (l *LoginController) Inject(
	responder *web.Responder,
	authManager *application.AuthManager,
	ph []LoginGetParameterHook,
	router web.ReverseRouter,
) {
	l.responder = responder
	l.authManager = authManager
	l.parameterHooks = ph
	l.router = router
}

// Get handler for logins (redirect)
func (l *LoginController) Get(c context.Context, request *web.Request) web.Result {
	redirecturl, ok := request.Params["redirecturl"]
	if !ok || redirecturl == "" {
		redirecturl = request.Request().Referer()
	}

	absolute, _ := l.router.Absolute(request, "", nil)
	if refURL, err := url.Parse(redirecturl); err != nil ||
		(refURL.Host != request.Request().Host && refURL.Host != absolute.Host) {
		redirecturl = absolute.String()
	}

	state := uuid.NewV4().String()
	l.authManager.StoreAuthState(request.Session(), state)
	request.Session().Store("auth.redirect", redirecturl)

	parameters := append([]oauth2.AuthCodeOption{}, l.authManager.AuthCodeOptions...)
	for _, hook := range l.parameterHooks {
		keyValue := hook.Parameters(c, request)
		for key, value := range keyValue {
			parameters = append(parameters, oauth2.SetAuthURLParam(key, value))
		}
	}

	redirectURL, _ := url.Parse(l.authManager.OAuth2Config(c, request).AuthCodeURL(state, parameters...))
	return l.responder.URLRedirect(redirectURL)
}
