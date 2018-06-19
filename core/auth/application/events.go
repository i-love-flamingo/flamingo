package application

import (
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// EventPublisher struct
	EventPublisher struct{}
)

// PublishLoginEvent dispatches the login event on the contexts event router
func (e *EventPublisher) PublishLoginEvent(ctx web.Context, event *domain.LoginEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(ctx, event)
}

// PublishLogoutEvent dispatches the logout event on the contexts event router
func (e *EventPublisher) PublishLogoutEvent(ctx web.Context, event *domain.LogoutEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(ctx, event)
}
