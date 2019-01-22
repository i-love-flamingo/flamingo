package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/v3/core/auth/application"
	"flamingo.me/flamingo/v3/core/auth/domain"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
)

type (
	// LogoutControllerInterface is the HTTP action provider implementation
	LogoutControllerInterface interface {
		Get(context.Context, *web.Request) web.Result
	}

	// LogoutController handles the logout
	LogoutController struct {
		responder      *web.Responder
		logger         flamingo.Logger
		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
		logoutRedirect LogoutRedirectAware
	}

	// DefaultLogoutRedirect helper
	DefaultLogoutRedirect struct {
		authManager *application.AuthManager
	}

	// LogoutRedirectAware to retrieve redirect urls
	LogoutRedirectAware interface {
		GetRedirectURL(context context.Context, u *url.URL) (*url.URL, error)
	}
)

var _ LogoutRedirectAware = new(DefaultLogoutRedirect)

// Inject DefaultLogoutRedirect dependencies
func (d *DefaultLogoutRedirect) Inject(manager *application.AuthManager) {
	d.authManager = manager
}

// GetRedirectURL builds default redirect URL for logout
func (d *DefaultLogoutRedirect) GetRedirectURL(c context.Context, u *url.URL) (*url.URL, error) {
	query := url.Values{}
	ru, _ := d.authManager.URL(c, "")
	query.Set("redirect_uri", ru.String())
	u.RawQuery = query.Encode()
	return u, nil
}

// Inject LogoutController dependencies
func (l *LogoutController) Inject(
	responder *web.Responder,
	logger flamingo.Logger,
	authManager *application.AuthManager,
	eventPublisher *application.EventPublisher,
	logoutRedirect LogoutRedirectAware,
) {
	l.responder = responder
	l.logger = logger
	l.authManager = authManager
	l.eventPublisher = eventPublisher
	l.logoutRedirect = logoutRedirect
}

// Get handler for logout
func (l *LogoutController) Get(ctx context.Context, request *web.Request) web.Result {
	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}

	ru, _ := l.authManager.URL(ctx, "")

	err := l.authManager.OpenIDProvider().Claims(&claims)
	if err != nil {
		l.logoutLocally(ctx, request)
		l.logger.Error("Logout locally only. Could not unmarshal raw fields", err.Error())
		return l.responder.URLRedirect(ru)
	}

	endURL, err := url.Parse(claims.EndSessionEndpoint)
	if err != nil {
		l.logoutLocally(ctx, request)
		l.logger.Error("Logout locally only. Could not parse end_session_endpoint claim to logout from IDP", err.Error())
		return l.responder.URLRedirect(ru)
	}

	redirectURL, redirectURLError := l.logoutRedirect.GetRedirectURL(ctx, endURL)
	if redirectURLError != nil {
		l.logoutLocally(ctx, request)
		l.logger.Error("Logout locally only. Could not fetch redirect URL for IDP logout", redirectURLError.Error())
		return l.responder.URLRedirect(ru)
	}

	l.logoutLocally(ctx, request)
	request.Session().AddFlash("successful logged out")

	return l.responder.URLRedirect(redirectURL)
}

func (l *LogoutController) logoutLocally(ctx context.Context, request *web.Request) {
	l.eventPublisher.PublishLogoutEvent(ctx, &domain.LogoutEvent{
		Session: request.Session(),
	})
	// todo fix
	// request.Session.G().Options.MaxAge = -1
}
