package application

import (
	"flamingo.me/flamingo/core/auth/domain"
	"flamingo.me/flamingo/framework/web"
)

type (
	// EventPublisher
	EventPublisher struct {
	}
)

func (e *EventPublisher) PublishLoginEvent(ctx web.Context, event *domain.LoginEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(ctx, event)
}

func (e *EventPublisher) PublishLogoutEvent(ctx web.Context, event *domain.LogoutEvent) {
	//publish to Flamingo default Event Router
	ctx.EventRouter().Dispatch(ctx, event)
}
