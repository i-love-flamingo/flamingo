package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	// LogoutController handles the logout
	LogoutController struct {
		responder.RedirectAware
		logger         flamingo.Logger
		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
		logoutRedirect LogoutRedirectAware
		myHost         string
	}

	// DefaultLogoutRedirect helper
	DefaultLogoutRedirect struct {
		authManager *application.AuthManager
		myHost      string
	}

	LogoutRedirectAware interface {
		GetRedirectUrl(context context.Context, u *url.URL) (string, error)
	}
)

var _ LogoutRedirectAware = new(DefaultLogoutRedirect)

// Inject DefaultLogoutRedirect dependencies
func (d *DefaultLogoutRedirect) Inject(manager *application.AuthManager, config *struct {
	MyHost string `inject:"config:auth.myhost"`
}) {
	d.authManager = manager
	d.myHost = config.MyHost
}

// GetRedirectUrl builds default redirect URL for logout
func (d *DefaultLogoutRedirect) GetRedirectUrl(_ context.Context, u *url.URL) (string, error) {
	query := url.Values{}
	query.Set("redirect_uri", d.myHost)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// Logout locally
func logout(r *web.Request) {
	delete(r.Session().Values, application.KeyAuthstate)
	delete(r.Session().Values, application.KeyToken)
	delete(r.Session().Values, application.KeyRawIDToken)
}

// Inject LogoutController dependencies
func (l *LogoutController) Inject(
	redirectAware responder.RedirectAware,
	logger flamingo.Logger,
	authManager *application.AuthManager,
	eventPublisher *application.EventPublisher,
	logoutRedirect LogoutRedirectAware,
	config *struct {
		MyHost string `inject:"config:auth.myhost"`
	},
) {
	l.RedirectAware = redirectAware
	l.logger = logger
	l.authManager = authManager
	l.eventPublisher = eventPublisher
	l.logoutRedirect = logoutRedirect
	l.myHost = config.MyHost
}

// Get handler for logout
func (l *LogoutController) Get(c context.Context, request *web.Request) web.Response {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}

	l.authManager.OpenIDProvider().Claims(&claims)
	endURL, parseError := url.Parse(claims.EndSessionEndpoint)
	if parseError != nil {
		logout(request)
		l.logger.Error("Logout locally only. Could not parse end_session_endpoint claim to logout from IDP", parseError.Error())
		return l.RedirectURL(l.myHost)
	}

	redirectURL, redirectURLError := l.logoutRedirect.GetRedirectUrl(web.ToContext(c), endURL)
	if redirectURLError != nil {
		logout(request)
		l.logger.Error("Logout locally only. Could not fetch redirect URL for IDP logout", redirectURLError.Error())
		return l.RedirectURL(l.myHost)
	}

	logout(request)

	request.Session().AddFlash("successful logged out", "warning")
	l.eventPublisher.PublishLogoutEvent(web.ToContext(c), &domain.LogoutEvent{
		Context: web.ToContext(c),
	})

	return l.RedirectURL(redirectURL)
}
