package middleware

import (
	"context"
	"net/url"

	"github.com/pkg/errors"

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
	}
)

func (m *SecurityMiddleware) Inject(r *web.Responder, s application.SecurityService, u RedirectUrlMaker, l flamingo.Logger, cfg *struct {
	LoginPathHandler              string `inject:"config:security.loginPath.handler"`
	LoginPathRedirectStrategy     string `inject:"config:security.loginPath.redirectStrategy"`
	LoginPathRedirectPath         string `inject:"config:security.loginPath.redirectPath"`
	AuthenticatedHomepageStrategy string `inject:"config:security.authenticatedHomepage.strategy"`
	AuthenticatedHomepagePath     string `inject:"config:security.authenticatedHomepage.path"`
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
}

func (m *SecurityMiddleware) HandleIfLoggedIn(action router.Action) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		if !m.securityService.IsLoggedIn(ctx, req.Session()) {
			redirectUrl := m.redirectUrl(ctx, req, m.loginPathRedirectStrategy, m.loginPathRedirectPath, req.Request().URL.String())
			return m.responder.RouteRedirect("auth.login", map[string]string{
				"redirecturl": redirectUrl.String(),
			})
		}
		return action(ctx, req)
	}
}

func (m *SecurityMiddleware) HandleIfLoggedOut(action router.Action) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		if !m.securityService.IsLoggedOut(ctx, req.Session()) {
			redirectUrl := m.redirectUrl(ctx, req, m.authenticatedHomepageStrategy, m.authenticatedHomepagePath, req.Request().Header.Get("Referer"))
			return m.responder.URLRedirect(redirectUrl)
		}
		return action(ctx, req)
	}
}

func (m *SecurityMiddleware) HandleIfGranted(action router.Action, permission string) router.Action {
	return func(ctx context.Context, req *web.Request) web.Response {
		if !m.securityService.IsGranted(ctx, req.Session(), permission, nil) {
			return m.responder.Forbidden(errors.Errorf("Missing permission %s for path %s.", permission, req.Request().URL.Path))
		}
		return action(ctx, req)
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
