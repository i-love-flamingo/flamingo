package middleware

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

	"fmt"
	"time"

	"strings"

	"flamingo.me/flamingo/core/security/application"
	"flamingo.me/flamingo/framework/flamingo"
	"flamingo.me/flamingo/framework/router"
	"flamingo.me/flamingo/framework/web"
)

const (
	ReferrerRedirectStrategy = "referrer"
	PathRedirectStrategy     = "path"
)

type (
	RedirectUrlMaker interface {
		URL(context.Context, string) (*url.URL, error)
	}

	SecurityMiddleware struct {
		responder        *web.Responder
		securityService  application.SecurityService
		redirectUrlMaker RedirectUrlMaker
		logger           flamingo.Logger

		loginPathHandler              string
		loginPathRedirectStrategy     string
		loginPathRedirectPath         string
		authenticatedHomepageStrategy string
		authenticatedHomepagePath     string
		eventLogging                  bool
	}
)

func (m *SecurityMiddleware) Inject(r *web.Responder, s application.SecurityService, u RedirectUrlMaker, l flamingo.Logger, cfg *struct {
	LoginPathHandler              string `inject:"config:security.loginPath.handler"`
	LoginPathRedirectStrategy     string `inject:"config:security.loginPath.redirectStrategy"`
	LoginPathRedirectPath         string `inject:"config:security.loginPath.redirectPath"`
	AuthenticatedHomepageStrategy string `inject:"config:security.authenticatedHomepage.strategy"`
	AuthenticatedHomepagePath     string `inject:"config:security.authenticatedHomepage.path"`
	EventLogging                  bool   `inject:"config:security.eventLogging"`
}) {
	m.responder = r
	m.securityService = s
	m.redirectUrlMaker = u
	m.logger = l
	m.loginPathHandler = cfg.LoginPathHandler
	m.loginPathRedirectStrategy = cfg.LoginPathRedirectStrategy
	m.loginPathRedirectPath = cfg.LoginPathRedirectPath
	m.authenticatedHomepageStrategy = cfg.AuthenticatedHomepageStrategy
	m.authenticatedHomepagePath = cfg.AuthenticatedHomepagePath
	m.eventLogging = cfg.EventLogging
}

func (m *SecurityMiddleware) HandleIfLoggedIn(action router.Action) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		if !m.securityService.IsLoggedIn(ctx, req.Session()) {
			return m.RedirectToLoginFallback(ctx, req)
		}
		return action(ctx, req)
	}
}

func (m *SecurityMiddleware) HandleIfLoggedOut(action router.Action) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		if !m.securityService.IsLoggedOut(ctx, req.Session()) {
			return m.RedirectToLogoutFallback(ctx, req)
		}
		return action(ctx, req)
	}
}

func (m *SecurityMiddleware) HandleIfGranted(action router.Action, permission string) router.Action {
	return m.handleForPermissionAndFallback(action, m.forbiddenAction(permission), true, permission)
}

func (m *SecurityMiddleware) HandleIfNotGranted(action router.Action, permission string) router.Action {
	return m.handleForPermissionAndFallback(action, m.forbiddenAction(permission), false, permission)
}

func (m *SecurityMiddleware) HandleIfGrantedWithFallback(action router.Action, fallback router.Action, permission string) router.Action {
	return m.handleForPermissionAndFallback(action, fallback, true, permission)
}

func (m *SecurityMiddleware) HandleIfNotGrantedWithFallback(action router.Action, fallback router.Action, permission string) router.Action {
	return m.handleForPermissionAndFallback(action, fallback, false, permission)
}

func (m *SecurityMiddleware) RedirectToLoginFallback(ctx context.Context, req *web.Request) web.Response {
	m.logIfNeeded(req, "request to only-authenticated page as unauthenticated user")
	redirectUrl := m.redirectUrl(ctx, req, m.loginPathRedirectStrategy, m.loginPathRedirectPath, req.Request().URL.String())
	return m.responder.RouteRedirect("auth.login", map[string]string{
		"redirecturl": redirectUrl.String(),
	})
}

func (m *SecurityMiddleware) RedirectToLogoutFallback(ctx context.Context, req *web.Request) web.Response {
	m.logIfNeeded(req, "request to only-unauthenticated page as authenticated user")
	redirectUrl := m.redirectUrl(ctx, req, m.authenticatedHomepageStrategy, m.authenticatedHomepagePath, req.Request().Header.Get("Referer"))
	return m.responder.URLRedirect(redirectUrl)
}

func (m *SecurityMiddleware) handleForPermissionAndFallback(action router.Action, fallback router.Action, ifGranted bool, permission string) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
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

func (m *SecurityMiddleware) forbiddenAction(permission string) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		return m.responder.Forbidden(errors.Errorf("Permission %s for path %s.", permission, req.Request().URL.Path))
	}
}

func (m *SecurityMiddleware) redirectUrl(ctx context.Context, req *web.Request, strategy string, path string, referrer string) *url.URL {
	var err error
	var generated *url.URL

	if strategy == ReferrerRedirectStrategy {
		generated, err = m.redirectUrlMaker.URL(ctx, referrer)
	} else if strategy == PathRedirectStrategy {
		generated, err = m.redirectUrlMaker.URL(ctx, path)
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
