package fake

import (
	"context"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/application/fake"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	LogoutController struct {
		responder.RedirectAware

		authManager *application.AuthManager
	}
)

func (l *LogoutController) Inject(
	redirectAware responder.RedirectAware,
	authManager *application.AuthManager,
) {
	l.RedirectAware = redirectAware
	l.authManager = authManager
}

func (l *LogoutController) Get(ctx context.Context, request *web.Request) web.Response {
	request.Session().Delete(fake.UserSessionKey)

	redirectUrl, _ := l.authManager.URL(ctx, "")

	return l.RedirectURL(redirectUrl.String())
}
