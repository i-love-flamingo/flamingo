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
	LogoutControllerInterface interface {
		Get(context.Context, *web.Request) web.Response
	}

	// LogoutController handles the logout
	LogoutController struct {
		responder.RedirectAware
		logger         flamingo.Logger
		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
		logoutRedirect LogoutRedirectAware
	}

	// DefaultLogoutRedirect helper
	DefaultLogoutRedirect struct {
		authManager *application.AuthManager
	}

	LogoutRedirectAware interface {
		GetRedirectURL(context context.Context, u *url.URL) (string, error)
	}
)

var _ LogoutRedirectAware = new(DefaultLogoutRedirect)

// Inject DefaultLogoutRedirect dependencies
func (d *DefaultLogoutRedirect) Inject(manager *application.AuthManager) {
	d.authManager = manager
}

// GetRedirectURL builds default redirect URL for logout
func (d *DefaultLogoutRedirect) GetRedirectURL(c context.Context, u *url.URL) (string, error) {
	query := url.Values{}
	ru, _ := d.authManager.URL(c, "")
	query.Set("redirect_uri", ru.String())
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// Logout locally
func logout(r *web.Request) {
	r.Session().Delete(application.KeyAuthstate)
	r.Session().Delete(application.KeyToken)
	r.Session().Delete(application.KeyRawIDToken)
	r.Session().Delete(application.KeyTokenExtras)

	// kill session
	r.Session().G().Options.MaxAge = -1
}

// Inject LogoutController dependencies
func (l *LogoutController) Inject(
	redirectAware responder.RedirectAware,
	logger flamingo.Logger,
	authManager *application.AuthManager,
	eventPublisher *application.EventPublisher,
	logoutRedirect LogoutRedirectAware,
) {
	l.RedirectAware = redirectAware
	l.logger = logger
	l.authManager = authManager
	l.eventPublisher = eventPublisher
	l.logoutRedirect = logoutRedirect
}

// Get handler for logout
func (l *LogoutController) Get(c context.Context, request *web.Request) web.Response {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}

	ru, _ := l.authManager.URL(c, "")

	l.authManager.OpenIDProvider().Claims(&claims)
	endURL, parseError := url.Parse(claims.EndSessionEndpoint)
	if parseError != nil {
		logout(request)
		l.logger.Error("Logout locally only. Could not parse end_session_endpoint claim to logout from IDP", parseError.Error())
		return l.RedirectURL(ru.String())
	}

	redirectURL, redirectURLError := l.logoutRedirect.GetRedirectURL(c, endURL)
	if redirectURLError != nil {
		logout(request)
		l.logger.Error("Logout locally only. Could not fetch redirect URL for IDP logout", redirectURLError.Error())
		return l.RedirectURL(ru.String())
	}

	logout(request)

	request.Session().AddFlash("successful logged out", "warning")
	l.eventPublisher.PublishLogoutEvent(c, &domain.LogoutEvent{
		Session: request.Session().G(),
	})

	return l.RedirectURL(redirectURL)
}
