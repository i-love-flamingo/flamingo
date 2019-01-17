package fake

import (
	"context"

	"flamingo.me/flamingo/core/auth/application"
	"flamingo.me/flamingo/core/auth/application/fake"
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
	"flamingo.me/flamingo/framework/web/responder"
)

type (
	LogoutController struct {
		responder.RedirectAware

		authManager    *application.AuthManager
		eventPublisher *application.EventPublisher
	}
)

func (l *LogoutController) Inject(
	redirectAware responder.RedirectAware,
	authManager *application.AuthManager,
	eventPublisher *application.EventPublisher,
) {
	l.RedirectAware = redirectAware
	l.authManager = authManager
	l.eventPublisher = eventPublisher
}

func (l *LogoutController) Get(ctx context.Context, request *web.Request) web.Response {
	request.Session().Delete(fake.UserSessionKey)
	l.eventPublisher.PublishLogoutEvent(ctx, &domain.LogoutEvent{
		Session: request.Session().G(),
	})

	redirectUrl, _ := l.authManager.URL(ctx, "")

	return l.RedirectURL(redirectUrl.String())
}
