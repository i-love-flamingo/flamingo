package application

import (
	"context"

	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/event"
)

// EventPublisher struct
type EventPublisher struct {
	router event.Router
}

func (e *EventPublisher) Inject(router event.Router) {
	e.router = router
}

// PublishLoginEvent dispatches the login event on the contexts event router
func (e *EventPublisher) PublishLoginEvent(ctx context.Context, event *domain.LoginEvent) {
	e.router.Dispatch(ctx, event)
}

// PublishLogoutEvent dispatches the logout event on the contexts event router
func (e *EventPublisher) PublishLogoutEvent(ctx context.Context, event *domain.LogoutEvent) {
	e.router.Dispatch(ctx, event)
}
