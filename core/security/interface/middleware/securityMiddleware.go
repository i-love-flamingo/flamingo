package middleware

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"strings"
	"time"

	"flamingo.me/flamingo/v3/core/security/application"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"flamingo.me/flamingo/v3/framework/web"
	"github.com/pkg/errors"
)

const (
	// ReferrerRedirectStrategy strategy to redirect to the supplied referrer
	ReferrerRedirectStrategy = "referrer"
	//PathRedirectStrategy strategy to redirect to the supplied path
	PathRedirectStrategy = "path"
)

type (
	// RedirectURLMaker helper function
	RedirectURLMaker interface {
		URL(context.Context, string) (*url.URL, error)
	}

	// RedirectURLMakerImpl is actual implementation for RedirectURLMaker interface
	RedirectURLMakerImpl struct {
		router web.ReverseRouter
	}

	// SecurityMiddleware to be used to secure controllers/routes
	SecurityMiddleware struct {
		responder        *web.Responder
		securityService  application.SecurityService
		redirectURLMaker RedirectURLMaker
		logger           flamingo.Logger

		loginPathHandler              string
		loginPathRedirectStrategy     string
		loginPathRedirectPath         string
		authenticatedHomepageStrategy string
		authenticatedHomepagePath     string
		eventLogging                  bool
	}
)

var _ RedirectURLMaker = new(RedirectURLMakerImpl)

// Inject dependencies
func (r *RedirectURLMakerImpl) Inject(router web.ReverseRouter) {
	r.router = router
}

// URL generates absolute url depending on provided path
func (r *RedirectURLMakerImpl) URL(ctx context.Context, redirectPath string) (*url.URL, error) {
	req := web.RequestFromContext(ctx)
	u, err := r.router.Absolute(req, "", nil)
	if err != nil {
		return u, err
	}
	u.Path = path.Join(u.Path, redirectPath)
	return u, nil
}

// Inject dependencies
func (m *SecurityMiddleware) Inject(r *web.Responder, s application.SecurityService, u RedirectURLMaker, l flamingo.Logger, cfg *struct {
	LoginPathHandler              string `inject:"config:security.loginPath.handler"`
	LoginPathRedirectStrategy     string `inject:"config:security.loginPath.redirectStrategy"`
	LoginPathRedirectPath         string `inject:"config:security.loginPath.redirectPath"`
	AuthenticatedHomepageStrategy string `inject:"config:security.authenticatedHomepage.strategy"`
	AuthenticatedHomepagePath     string `inject:"config:security.authenticatedHomepage.path"`
	EventLogging                  bool   `inject:"config:security.eventLogging"`
}) {
	m.responder = r
	m.securityService = s
	m.redirectURLMaker = u
	m.logger = l
	m.loginPathHandler = cfg.LoginPathHandler
	m.loginPathRedirectStrategy = cfg.LoginPathRedirectStrategy
	m.loginPathRedirectPath = cfg.LoginPathRedirectPath
	m.authenticatedHomepageStrategy = cfg.AuthenticatedHomepageStrategy
	m.authenticatedHomepagePath = cfg.AuthenticatedHomepagePath
	m.eventLogging = cfg.EventLogging
}

// HandleIfLoggedIn allows a controller to be used for logged in users
func (m *SecurityMiddleware) HandleIfLoggedIn(action web.Action) web.Action {
	return func(ctx context.Context, req *web.Request) web.Result {
		if !m.securityService.IsLoggedIn(ctx, req.Session()) {
			return m.RedirectToLoginFallback(ctx, req)
		}
		return action(ctx, req)
	}
}

// HandleIfLoggedOut allows a controller to be used for logged out users
func (m *SecurityMiddleware) HandleIfLoggedOut(action web.Action) web.Action {
	return func(ctx context.Context, req *web.Request) web.Result {
		if !m.securityService.IsLoggedOut(ctx, req.Session()) {
			return m.RedirectToLogoutFallback(ctx, req)
		}
		return action(ctx, req)
	}
}

// HandleIfGranted allows a controller to be used with a given permission
func (m *SecurityMiddleware) HandleIfGranted(action web.Action, permission string) web.Action {
	return m.handleForPermissionAndFallback(action, m.forbiddenAction(permission), true, permission)
}

// HandleIfNotGranted allows a controller not to be used with a given permission
func (m *SecurityMiddleware) HandleIfNotGranted(action web.Action, permission string) web.Action {
	return m.handleForPermissionAndFallback(action, m.forbiddenAction(permission), false, permission)
}

// HandleIfGrantedWithFallback is HandleIfGranted with a fallback action
func (m *SecurityMiddleware) HandleIfGrantedWithFallback(action web.Action, fallback web.Action, permission string) web.Action {
	return m.handleForPermissionAndFallback(action, fallback, true, permission)
}

// HandleIfNotGrantedWithFallback is HandleIfNotGranted with a fallback action
func (m *SecurityMiddleware) HandleIfNotGrantedWithFallback(action web.Action, fallback web.Action, permission string) web.Action {
	return m.handleForPermissionAndFallback(action, fallback, false, permission)
}

// RedirectToLoginFallback fallback helper action
func (m *SecurityMiddleware) RedirectToLoginFallback(ctx context.Context, req *web.Request) web.Result {
	m.logIfNeeded(req, "request to only-authenticated page as unauthenticated user")
	redirectURL := m.redirectURL(ctx, req, m.loginPathRedirectStrategy, m.loginPathRedirectPath, req.Request().URL.String())
	return m.responder.RouteRedirect("auth.login", map[string]string{
		"redirecturl": redirectURL.String(),
	})
}

// RedirectToLogoutFallback fallback helper action
func (m *SecurityMiddleware) RedirectToLogoutFallback(ctx context.Context, req *web.Request) web.Result {
	m.logIfNeeded(req, "request to only-unauthenticated page as authenticated user")
	redirectURL := m.redirectURL(ctx, req, m.authenticatedHomepageStrategy, m.authenticatedHomepagePath, req.Request().Header.Get("Referer"))
	return m.responder.URLRedirect(redirectURL)
}

func (m *SecurityMiddleware) handleForPermissionAndFallback(action web.Action, fallback web.Action, ifGranted bool, permission string) web.Action {
	return func(ctx context.Context, req *web.Request) web.Result {
		granted := m.securityService.IsGranted(ctx, req.Session(), permission, nil)
		if (ifGranted && !granted) || (!ifGranted && granted) {
			explanation := "without"
			if !ifGranted {
				explanation = "with"
			}
			m.logIfNeeded(req, fmt.Sprintf("request to protected page %s permission %s", explanation, permission))
			return fallback(ctx, req)
		}
		return action(ctx, req)
	}
}

func (m *SecurityMiddleware) forbiddenAction(permission string) web.Action {
	return func(ctx context.Context, req *web.Request) web.Result {
		return m.responder.Forbidden(errors.Errorf("Permission %s for path %s.", permission, req.Request().URL.Path))
	}
}

func (m *SecurityMiddleware) redirectURL(ctx context.Context, req *web.Request, strategy string, path string, referrer string) *url.URL {
	var err error
	var generated *url.URL

	if strategy == ReferrerRedirectStrategy {
		generated, err = m.redirectURLMaker.URL(ctx, referrer)
	} else if strategy == PathRedirectStrategy {
		generated, err = m.redirectURLMaker.URL(ctx, path)
	}

	if err != nil {
		m.logger.Error(err)
	}

	return generated
}

func (m *SecurityMiddleware) logIfNeeded(r *web.Request, message string) {
	if m.eventLogging {
		m.logger.WithField("security", "middleware").
			WithField("Date", time.Now().Format(time.RFC3339)).
			WithField("Path", r.Request().URL.Path).
			WithField("RemoteAddress", strings.Join(r.RemoteAddress(), ", ")).
			Info(message)
	}
}
