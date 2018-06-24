package interfaces

import (
	"context"
	"net/url"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

type (
	// LoginController handles the login redirect
	LoginController struct {
		responder.RedirectAware
		authManager *application.AuthManager
		myHost      string
	}

	// LogoutController handles the logout
	LogoutController struct {
		responder.RedirectAware
		logger         flamingo.Logger
		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
		logoutRedirect LogoutRedirectAware
		myHost         string
	}

	// CallbackController handles the oauth2.0 callback
	CallbackController struct {
		responder.RedirectAware
		responder.ErrorAware
		authManager    *application.AuthManager
		logger         flamingo.Logger
		eventPublisher *application.EventPublisher
	}

	// DefaultLogoutRedirect helper
	DefaultLogoutRedirect struct {
		authManager *application.AuthManager
		myHost      string
	}

	LogoutRedirectAware interface {
		GetRedirectUrl(c web.Context, u *url.URL) (string, error)
	}
)

// Inject DefaultLogoutRedirect dependencies
func (d *DefaultLogoutRedirect) Inject(manager *application.AuthManager, config *struct {
	MyHost string `inject:"config:auth.myhost"`
}) {
	d.authManager = manager
	d.myHost = config.MyHost
}

// Build default redirect URL for logout
func (d *DefaultLogoutRedirect) GetRedirectUrl(c web.Context, u *url.URL) (string, error) {
	query := url.Values{}
	query.Set("redirect_uri", d.myHost)
	u.RawQuery = query.Encode()
	return u.String(), nil
}

// Logout locally
func logout(c web.Context) {
	delete(c.Session().Values, application.KeyAuthstate)
	delete(c.Session().Values, application.KeyToken)
	delete(c.Session().Values, application.KeyRawIDToken)
}

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
func (l *LoginController) Get(c_ context.Context, _ *web.Request) web.Response {
	c := c_.Value(web.CONTEXT).(web.Context)

	redirecturl, err := c.Param1("redirecturl")
	if err != nil || redirecturl == "" {
		redirecturl = c.Request().Referer()
	}

	if refURL, err := url.Parse(redirecturl); err != nil || refURL.Host != c.Request().Host {
		redirecturl = l.myHost
	}

	state := uuid.NewV4().String()
	c.Session().Values["auth.state"] = state
	c.Session().Values["auth.redirect"] = redirecturl

	return l.RedirectURL(l.authManager.OAuth2Config().AuthCodeURL(state))
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
func (l *LogoutController) Get(c_ context.Context, _ *web.Request) web.Response {
	c := c_.Value(web.CONTEXT).(web.Context)

	var claims struct {
		EndSessionEndpoint string `json:"end_session_endpoint"`
	}

	l.authManager.OpenIDProvider().Claims(&claims)
	endUrl, parseError := url.Parse(claims.EndSessionEndpoint)
	if parseError != nil {
		logout(c)
		l.logger.Error("Logout locally only. Could not parse end_session_endpoint claim to logout from IDP", parseError.Error())
		return l.RedirectURL(l.myHost)
	}

	redirectUrl, redirectUrlError := l.logoutRedirect.GetRedirectUrl(c, endUrl)
	if redirectUrlError != nil {
		logout(c)
		l.logger.Error("Logout locally only. Could not fetch redirect URL for IDP logout", redirectUrlError.Error())
		return l.RedirectURL(l.myHost)
	}

	logout(c)

	c.Session().AddFlash("successful logged out", "warning")
	l.eventPublisher.PublishLogoutEvent(c, &domain.LogoutEvent{
		Context: c,
	})

	return l.RedirectURL(redirectUrl)
}

// Inject CallbackController dependencies
func (cc *CallbackController) Inject(
	redirectAware responder.RedirectAware,
	errorAware responder.ErrorAware,
	authManager *application.AuthManager,
	logger flamingo.Logger,
	eventPublisher *application.EventPublisher,
) {
	cc.RedirectAware = redirectAware
	cc.ErrorAware = errorAware
	cc.authManager = authManager
	cc.logger = logger
	cc.eventPublisher = eventPublisher
}

// Get handler for callbacks
func (cc *CallbackController) Get(c_ context.Context, _ *web.Request) web.Response {
	c := c_.Value(web.CONTEXT).(web.Context)

	// Verify state and errors.
	defer delete(c.Session().Values, application.KeyAuthstate)

	if c.Session().Values[application.KeyAuthstate] != c.MustQuery1("state") {
		cc.logger.Error("Invalid State", c.Session().Values[application.KeyAuthstate], c.MustQuery1("state"))
		return cc.Error(c, errors.New("Invalid State"))
	}

	finish := c.Profile("auth.callback", "code: "+c.MustQuery1("code"))
	oauth2Token, err := cc.authManager.OAuth2Config().Exchange(c, c.MustQuery1("code"))
	finish()
	if err != nil {
		cc.logger.Error("core.auth.callback Error OAuth2Config Exchange", err)
		return cc.Error(c, errors.WithStack(err))
	}

	c.Session().Values[application.KeyToken] = oauth2Token
	c.Session().Values[application.KeyRawIDToken], err = cc.authManager.ExtractRawIDToken(oauth2Token)
	if err != nil {
		cc.logger.Error("core.auth.callback Error ExtractRawIDToken", err)
		return cc.Error(c, errors.WithStack(err))
	}
	cc.eventPublisher.PublishLoginEvent(c, &domain.LoginEvent{Context: c})
	cc.logger.Debug("successful logged in and saved tokens", oauth2Token)
	c.Session().AddFlash("successful logged in", "info")

	if redirect, ok := c.Session().Values["auth.redirect"]; ok {
		delete(c.Session().Values, "auth.redirect")
		return cc.RedirectURL(redirect.(string))
	}
	return cc.Redirect("home", nil)
}
